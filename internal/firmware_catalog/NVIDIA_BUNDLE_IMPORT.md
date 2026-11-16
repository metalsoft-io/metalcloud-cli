# NVIDIA HGX/DGX firmware bundle import

NVIDIA publishes no machine-readable catalog. Operators download HGX/DGX bundles from the NVIDIA
Enterprise Support Portal or NGC and describe them with a JSON manifest. The CLI turns the manifest
into a MetalSoft firmware catalog (one catalog per bundle version) through the normal catalog and
binary endpoints.

## Status

### Done in nodejs-metalcloud (inventory)

- `nvidia` catalog vendor (`ServerFirmwareCatalogVendor.SERVER_FIRMWARE_CATALOG_VENDOR_NVIDIA`).
- Firmware binaries have `isBundle`, `bundleSteps` and `sha256`. A bundle needs at least one step, and
  only bundles can have steps (400 otherwise). Step fields and the relations between steps (duplicate ids,
  unknown dependency, self dependency, cycles) are validated by the API too.
- `PUT /firmware/binary/:id` (PUT semantics: optional fields that are left out are reset).
- Unique key `(catalogId, externalId)` on binaries; duplicates return 409.
- Baseline search: a DGX (`vendor = NVIDIA`) is searched in NVIDIA catalogs only. On an OEM server, the
  components with a `baseboardId` are also searched in NVIDIA catalogs, and an NVIDIA bundle replaces the
  OEM result for those components. Platform = `"<server vendor> <server model>"`, compared ignoring case
  and all whitespace with the bundle `vendorSupportedSystems[].id` and prerequisite `platform`: `"NVIDIA DGX H100"`
  (NVIDIA's name) is the server `NVIDIA`/`DGXH100` (what its BMC reports). Hyphens and suffixes are compared as
  written.
- Migration SQL: `apps/inventory/migrations/add_nvidia_firmware_bundles.sql` (to be moved to the schema project).

### Done in this repo

- `nvidia_manifest.go`: `LoadNvidiaManifest` / `ParseNvidiaManifest` (rejects unknown fields),
  `Validate`, `VerifyNvidiaPackageFiles` (exists, regular file, inside the folder including through
  symlinks, SHA-256), `PackageFiles`, `PackageFileSha256`, `Warnings`. Every issue about a step starts with
  `Step "<id>":`, every issue about a prerequisite with `prerequisites[<n>]:`, and all issues are reported at once.
- `vendor_catalog_nvidia.go`: `processNvidiaCatalog` builds the catalog and its binaries from the manifest (see
  [Catalog and binaries](#catalog-and-binaries)); `nvidia` is in `ValidVendors`.
- `CreateMetalsoftCatalog` creates or updates NVIDIA catalogs (see [Re-importing](#re-importing-a-bundle)).
- `firmware-catalog create --vendor nvidia` and `--force`.
- `firmware-binary list --filter-is-bundle/--filter-component-class`, `firmware-binary get` (lists the bundle
  steps and prerequisites), `firmware-binary update` (changes fields of a binary, e.g. its prerequisites) and
  `firmware-binary config-example --bundle`.
- Tests: `nvidia_manifest_test.go` and `vendor_catalog_nvidia_test.go`, with the fixture in
  `testdata/nvidia/hgx-h100/` (dummy `.fwpkg` files with known checksums). `.gitignore` ignores `*.json`, so an
  exception was added for `testdata/**/*.json`.

## Manifest format

See `testdata/nvidia/hgx-h100/manifest.json` for a complete example.

| Field | Required | Notes |
|---|---|---|
| `schemaVersion` | yes | `1` |
| `bundleName`, `bundleVersion` | yes | Catalog `vendorId` becomes `<bundleName>@<bundleVersion>` |
| `bundleFile` | yes | The main package; must be the `packageFile` of at least one step |
| `targetPlatforms` | yes | `"<server vendor> <server model>"` as stored on the server, e.g. `"NVIDIA DGXH100"` (or `"NVIDIA DGX H100"`: case and spaces are ignored), `"Supermicro SYS-821GE-TNHR"` |
| `releaseDate` | yes | `YYYY-MM-DD` |
| `releaseNotesUrl` | no | http(s) |
| `updateSeverity` | no | `critical`, `recommended`, `optional`, `unknown` (default `unknown`) |
| `downloadUrls` | no | Map of package file to a vendor URL. Files without one use the repository URL |
| `steps[]` | yes | `id`, `componentClass`, `targetPattern`, `packageFile`, `packageVersion`, `sha256`, `activation`, `dependsOnStep` (step id or `null`), `verifyVersion` (optional) |
| `prerequisites[]` | no | `platform` (optional, one of `targetPlatforms`), `componentClass`, `componentModel`, `componentName`, `minimumVersion`. Copied to the bundle binary; the bundle is not scheduled on a server of that platform until the component has at least that version. The usual entries are the minimum host BIOS and BMC from the NVIDIA release notes |

`componentClass` must be a `ServerComponentClass` value (mirrored in `ValidServerComponentClasses`, keep it
in sync with `libs/dto/server.component.dto.ts`). `activation` is `none`, `host_reset`, `bmc_reset`
(the host BMC), `hgx_reset` (the HGX HMC, Manager.Reset on the baseboard's manager) or `chassis_power_cycle`.
`verifyVersion` is the version the component reports once the step is active (a bundle's `packageVersion` is the
bundle release, not what a GPU reports); when set, the server microservice fails the update if the BMC reports
another version after activation.

The inventory versions a bundle per component by the step that updates it: `verifyVersion`, else the step's
`packageVersion`. That version is what the audit, the available versions (and so the target version), the baseline
minimum versions and the search's newest-binary choice compare with the installed version; the bundle release is
kept beside it as `bundleVersion`. Give every step the version its components report, in their own numbering
(`96.00.A5.00.01` for a GPU VBIOS, not `25.06.1`): `NvidiaBundleManifest.Warnings()` flags a step without a
`verifyVersion` whose `packageVersion` is the `bundleVersion`. `targetPattern` is checked here with Go's RE2, and the inventory matches it as a
case-insensitive JavaScript regular expression against the Redfish member id of the component: keep
patterns to the syntax both share (anchors, groups, alternation, character classes, `\d`, quantifiers),
with no lookarounds and no inline flags.

### How a bundle is applied: order and grouping

A bundle is applied as one job group with one job per **package delivery**, run in sequence:

- Components of the bundle that share a `packageFile` (and the same Redfish target, `HGX_0` on a DGX) are one job:
  the package is pushed once, the BMC task is waited for, each distinct `activation` of those steps runs once, then
  every component's version is checked. A BMC that cannot target single components inside a package (the DGX pulls
  large packages with `SimpleUpdate`, which takes no `Targets`) flashes the whole package anyway, and its HMC orders
  the components inside it.
- Jobs run in the order of the manifest when a step declares `dependsOnStep` (a step after the step it depends on,
  otherwise in the listed order). Without any `dependsOnStep`, the default HGX order applies, by component class:
  `gpu_erot`, `hgx_hmc`, `hgx_fpga`, `pcie_switch`, `pcie_retimer`, `nvswitch`, `gpu`; other classes after them.
  A job runs at the position of its earliest step.
- The activations of one job run in this order: `bmc_reset`, `hgx_reset`, `host_reset`, `chassis_power_cycle`;
  `host_reset` is skipped when a `chassis_power_cycle` follows.
- A step may not depend on a step that ships in a package delivered by a later job (the upgrade is refused). A
  bundle upgrade staged until the next reboot is refused when a step depends on another job's step that needs an
  activation, since staging cannot make that step active first.
- When a job fails, the remaining jobs do not run and the server stays in `updating_firmware`. Retrying the job
  continues the bundle from there; skipping it abandons the rest of the bundle, and the finish job puts the server
  back in its previous status. When the upgrade finishes, the firmware inventory is read again, so components a
  package updated without being scheduled show their new version.

In `testdata/nvidia/hgx-h100/manifest.json`, `erot`, `hmc` and `gpu` ship in `nvfw_HGX_H100x8_250612.fwpkg` and
are one job (activations `hgx_reset`, then `chassis_power_cycle`); `retimer` has its own package and runs second,
after `hmc` (`host_reset`).

### Delivery to the HGX tray

A step whose components sit on a GPU baseboard (`HGX_0` on a DGX) is only delivered with its Redfish `Targets`,
following NVIDIA's DGX procedure (multipart upload to `UpdateService/upload` with
`UpdateParameters: {"Targets": ["/redfish/v1/UpdateService/FirmwareInventory/HGX_0"]}`):

- pushed multipart, or with `HttpPushUri` when the BMC has `HttpPushUriTargets`, also when the image is above the
  pull threshold;
- pulled with `SimpleUpdate` only when its ActionInfo lists `Targets` (the DGX's does not);
- otherwise refused before anything is downloaded, with the image size, `MaxImageSizeBytes` and the methods the BMC
  offers in the message. On a DGX this happens for a package larger than its `MaxImageSizeBytes` (402 MiB on the
  captured DGX H100).

After the step's activation, the tray members are listed again in the BMC's `FirmwareInventory` before their
versions are read (the server microservice waits up to `firmwareHgxRecoveryTimeoutMinutes`, 20 by default).

### nvfwupd config keys and MetalCloud

NVIDIA's `nvfwupd` takes a platform-agnostic YAML config. How its keys relate to a MetalCloud bundle (nothing below
is read from an nvfwupd file; the manifest has no extra fields for them yet):

| nvfwupd key | MetalCloud |
|---|---|
| `TargetPlatform` (`DGX`, `HGX`, …) | The bundle's `targetPlatforms` (`"<vendor> <model>"` as stored on the server) |
| `FWUpdateFilePath` | A step's `packageFile`, next to the bundle in the firmware repository |
| `FwUpdateMethod: MultipartHttpPushUri` / `HttpPushUri` | Chosen by the server microservice from the BMC's UpdateService (multipart first); a tray step is never pulled without `Targets` |
| `MultipartHttpPushUri`, `HttpPushUri`, `TaskServiceUri` overrides | Read from the BMC's UpdateService and from the task the BMC returns |
| `UpdateParametersTargets` | The component's update target from the inventory: `…/FirmwareInventory/HGX_0` for every tray member |
| `MultipartOptions: ForceUpdate` | Not set: the BMC keeps its default (on a DGX only newer versions are flashed) |
| `SpecialUpdateParameters: ApplyTime` | `@Redfish.OperationApplyTime`, from the step `activation` and the staged option |
| `BMCResetParameters` | A step `activation` (`bmc_reset`, `hgx_reset`). `ResetToDefaultsType: ResetAll` is a reset to factory defaults and is never sent |
| `BMC_IP`, `RF_USERNAME`, `RF_PASSWORD`, `Targets` (systems) | The server's BMC address and credentials; one upgrade per server |

## Importing a bundle

```
metalcloud-cli firmware-catalog create --name "NVIDIA HGX H100 25.06.1" --vendor nvidia --update-type offline \
  --vendor-local-catalog-path ./h100/manifest.json --vendor-local-binaries-path ./h100 \
  --repo-base-url https://repo.example.com/nvidia --upload-binaries --repo-ssh-host repo.example.com:22 --repo-ssh-user upload
```

`--vendor-local-catalog-path` is the manifest and `--vendor-local-binaries-path` the folder with the package files.
`--vendor-url` and `--download-binaries` are refused for NVIDIA. `--repo-base-url` is required unless every package
file has a `downloadUrls` entry. `--server-types` is ignored with a warning: the supported systems come from
`targetPlatforms`, and a server-type lookup matches by vendor, so it would never find HGX trays in OEM chassis.

### Catalog and binaries

| Catalog field | Value |
|---|---|
| `VendorId` | `<bundleName>@<bundleVersion>` |
| `VendorReleaseTimestamp` | `releaseDate` |
| `VendorServerTypesSupported` | `targetPlatforms` |
| `VendorConfiguration` | `{schemaVersion, bundleName, bundleVersion, releaseNotesUrl}` |

The bundle binary comes from `bundleFile`: `isBundle`, every step in `bundleSteps`, `sha256`, the
`prerequisites`, one supported device per step (`{id: componentClass, model: targetPattern}`), one supported system per
target platform, and `rebootRequired` when a step's activation is `host_reset` or `chassis_power_cycle`. Each other
package file of the steps becomes a binary of its own (`isBundle` false), named after the steps that use it. A
binary's `externalId` is its package file: the bundle steps refer to it by `packageFile`, and it is the upload path.

### Re-importing a bundle

Importing the same `bundleName@bundleVersion` again updates the existing catalog and its binaries (matched by
`externalId`) instead of creating new ones, so the ids, and the baselines that use them, are kept. If a package file
already in the catalog has another checksum, the import stops and names the files; `--force` replaces them. The
other vendors keep the create-only behaviour.

### Checking the result

- `metalcloud-cli firmware-binary list --filter-catalog-id <id> --filter-is-bundle true`, then
  `firmware-binary get <id>` for the steps and prerequisites.
- `metalcloud-cli server firmware search-dry-run <server_id> --catalog-ids <id>` shows, per component, the binaries
  that matched (with their rules) and the ones rejected (with the reason).
- `metalcloud-cli server firmware generate-audit --include-preflight` runs the safety checks an upgrade would run
  now and lists the jobs it would create, in order (one per bundle step).
