package firmware_catalog

import (
	"context"
	"os"
	"strings"
	"testing"
)

const mockCatalogXML = `<?xml version="1.0" encoding="utf-16"?>
<Manifest baseLocation="downloads.dell.com" baseLocationAccessProtocols="HTTPS" dateTime="2025-04-22T01:46:12+05:30" identifier="48912fae-2b46-4b4c-bafe-2c709c7b0ad2" releaseID="RHDC6" version="25.04.28" predecessorID="02c9dbe4-de0f-495c-8e1d-8b804b0894e6">
  <ReleaseNotes>
    <Display lang="en">Release Notes</Display>
  </ReleaseNotes>
  <InventoryComponent dateTime="2025-03-13T04:09:54Z" dellVersion="A00" hashMD5="98b8a4f6bb2bf57722276c409f80fb2a" osCode="LIN64" path="FOLDER12838322M/3/invcol_LN64_4J9X7_25_03_00_314_A00.BIN" releaseDate="March 13, 2025" releaseID="4J9X7" schemaVersion="2.0" vendorVersion="25.03.00" />
  <InventoryComponent dateTime="2025-03-13T04:09:54Z" dellVersion="A00" hashMD5="7b4659e8a2817ba076b223bdd6cd6462" osCode="WIN64" path="FOLDER12838282M/3/invCol_WIN64_4J9X7_25_03_00_314_A00.exe" releaseDate="March 13, 2025" releaseID="4J9X7" schemaVersion="2.0" vendorVersion="25.03.00" />
  <SoftwareComponent dateTime="2017-03-24T21:05:18+05:30" dellVersion="A00" hashMD5="3979c65df3c67a5342d707af89923de5" packageID="H09VC" packageType="LLXP" path="FOLDER04177723M/1/Serial-ATA_Firmware_H09VC_LN_MA8F_A00.BIN" rebootRequired="false" releaseDate="March 24, 2017" releaseID="H09VC" schemaVersion="2.4" size="40821221" vendorVersion="MA8F">
    <Name>
      <Display lang="en"><![CDATA[Seagate MA8F for model number(s) ST6000NM0024-1US17Z..]]></Display>
    </Name>
    <ComponentType value="FRMW">
      <Display lang="en"><![CDATA[Firmware]]></Display>
    </ComponentType>
    <Description>
      <Display lang="en"><![CDATA[This release contains firmware version MA8F for Seagate drives. Vendor model numbers ST6000NM0024-1US17Z..]]></Display>
    </Description>
    <LUCategory value="Serial ATA">
      <Display lang="en"><![CDATA[Serial ATA]]></Display>
    </LUCategory>
    <Category value="SA">
      <Display lang="en"><![CDATA[Serial ATA]]></Display>
    </Category>
    <ImportantInfo URL="https://www.dell.com/support/home/en-us/drivers/DriversDetails?driverId=H09VC">
      <Display lang="en"><![CDATA[NA]]></Display>
    </ImportantInfo>
    <SupportedDevices>
      <Device componentID="103733" embedded="0">
        <Display lang="en"><![CDATA[Makara SATA 512e]]></Display>
      </Device>
    </SupportedDevices>
    <RevisionHistory>
      <Display lang="en"><![CDATA[This release contains firmware version MA8F for Seagate drives. Vendor model numbers ST6000NM0024-1US17Z..]]></Display>
    </RevisionHistory>
    <Criticality value="1">
      <Display lang="en"><![CDATA[Recommended]]></Display>
    </Criticality>
    <SupportedSystems>
      <Brand key="3" prefix="PE">
        <Display lang="en"><![CDATA[PowerEdge]]></Display>
        <Model systemID="0627" systemIDType="BIOS">
          <Display lang="en"><![CDATA[R730xd]]></Display>
        </Model>
        <Model systemID="0639" systemIDType="BIOS">
          <Display lang="en"><![CDATA[R430]]></Display>
        </Model>
        <Model systemID="063A" systemIDType="BIOS">
          <Display lang="en"><![CDATA[R530]]></Display>
        </Model>
        <Model systemID="063B" systemIDType="BIOS">
          <Display lang="en"><![CDATA[T430]]></Display>
        </Model>
        <Model systemID="06A6" systemIDType="BIOS">
          <Display lang="en"><![CDATA[R330]]></Display>
        </Model>
        <Model systemID="0600" systemIDType="BIOS">
          <Display lang="en"><![CDATA[R730]]></Display>
        </Model>
        <Model systemID="06A7" systemIDType="BIOS">
          <Display lang="en"><![CDATA[T330]]></Display>
        </Model>
        <Model systemID="0602" systemIDType="BIOS">
          <Display lang="en"><![CDATA[T630]]></Display>
        </Model>
      </Brand>
    </SupportedSystems>
  </SoftwareComponent>
  <SoftwareComponent dateTime="2024-10-15T05:03:16-05:00" dellVersion="A00" hashMD5="458fff12211093ea98bdf319d194da0e" packageID="YN8W5" packageType="LLXP" path="FOLDER12213519M/1/Systems-Management_Application_YN8W5_LN64_5.4.0.0_A00.BIN" rebootRequired="false" releaseDate="October 15, 2024" releaseID="YN8W5" schemaVersion="1.0" size="55908725" vendorVersion="5.4.0.0">
    <Name>
      <Display lang="en"><![CDATA[Dell iDRAC Service Module Embedded Package v5.4.0.0, A00]]></Display>
    </Name>
    <ComponentType value="APAC">
      <Display lang="en"><![CDATA[Application]]></Display>
    </ComponentType>
    <Description>
      <Display lang="en"><![CDATA[iDRAC Service Module (iSM) is a lightweight software service that better integrates operating system (OS) features with iDRAC and can be installed on Dell’s 14G or later generation of PowerEdge servers. iSM provides OS-related information to the iDRAC and adds capabilities such as LC log event replication into the OS log, Metric injection from operating system to iDRAC telemetry,WMI support (including storage), iDRAC SNMP alerts via OS, iDRAC hard reset and remote full Power Cycle. iSM automates SupportAssist report collection process for iDRAC leading to faster issue resolution. iSM has very little impact on the host processor and smaller memory footprint than in-band agents such as Dell OpenManage Server Administrator (OMSA), thus expanding iDRAC management into supported host operating systems.]]></Display>
    </Description>
    <LUCategory value="Systems Management">
      <Display lang="en"><![CDATA[Systems Management]]></Display>
    </LUCategory>
    <Category value="SM">
      <Display lang="en"><![CDATA[Systems Management]]></Display>
    </Category>
    <SupportedDevices>
      <Device componentID="104684" embedded="1">
        <Display lang="en"><![CDATA[iSM LC DUP]]></Display>
        <RollbackInformation fmpWrapperIdentifier="4ABC2C44-7F05-447F-B139-B66D2E72E7CF" fmpWrapperVersion="1.0" impactsTPMmeasurements="true" rollbackIdentifier="EE59FC83-783B-452B-BD14-B85D5600D4A9" rollbackTimeout="1200" rollbackVolume="MAS028" />
        <PayloadConfiguration>
          <Image filename="OM-iSM-Dell-Web-LX-5.4.0.0.tar.gz" id="5DD5A8BA-1958-4673-BE77-40B69680AF5D" skip="false" type="APAC" version="5.4.0.0" />
          <Image filename="OM-iSM-Dell-Web-LX-5.4.0.0.tar.gz.sign" id="E166C545-82A9-4D5D-8493-B834850F9C7A" skip="false" type="APAC" version="5.4.0.0" />
          <Image filename="OM-iSM-Dell-Web-X64-5.4.0.0.exe" id="5015744F-F938-40A8-B695-5456E9055504" skip="false" type="APAC" version="5.4.0.0" />
          <Image filename="ISM-Dell-Web-5.4.0.0-VIB-ESX8i-Live.zip" id="305D3492-F3B0-11EC-B939-0242AC120002" skip="false" type="APAC" version="5.4.0.0" />
          <Image filename="ISM-Dell-Web-5.4.0.0-VIB-ESX7i-Live.zip" id="1161DCFA-DB42-4EF1-A2B2-2D7D17091634" skip="false" type="APAC" version="5.4.0.0" />
          <Image filename="RPM-GPG-KEY-dell" id="0538B4E9-DA4D-402A-9D96-A4A55EE2234C" skip="false" type="APAC" version="" />
          <Image filename="sha256sum" id="06F61B54-58E2-41FB-8CE3-B7137A60E4B7" skip="false" type="APAC" version="" />
        </PayloadConfiguration>
      </Device>
    </SupportedDevices>
    <SupportedSystems>
      <Brand key="3" prefix="PE">
        <Display lang="en"><![CDATA[Enterprise Servers]]></Display>
        <Model systemID="0B1B" systemIDType="BIOS">
          <Display lang="en"><![CDATA[XR5610]]></Display>
        </Model>
      </Brand>
    </SupportedSystems>
    <ImportantInfo URL="https://www.dell.com/support">
      <Display lang="en"><![CDATA[NA]]></Display>
    </ImportantInfo>
    <Criticality value="1">
      <Display lang="en"><![CDATA[Recommended]]></Display>
    </Criticality>
    <FMPWrappers>
      <FMPWrapperInformation digitalSignature="false" filePathName="ApplicationWrapper.efi" identifier="4ABC2C44-7F05-447F-B139-B66D2E72E7CF" name="LC-ISM">
        <Inventory source="LCL" supported="true" />
        <Update rollback="false" supported="true" />
      </FMPWrapperInformation>
    </FMPWrappers>
  </SoftwareComponent>
  <SoftwareComponent dateTime="2024-05-27T08:45:23+05:30" dellVersion="A00" hashMD5="83ad12b7eab4c50099dccffdb1066829" packageID="5JJFF" packageType="LW64" path="FOLDER11615172M/1/Express-Flash-PCIe-SSD_Firmware_5JJFF_WN64_2.0.0_A00.EXE" rebootRequired="false" releaseDate="May 27, 2024" releaseID="5JJFF" schemaVersion="2.4" size="18553112" vendorVersion="2.0.0">
    <Name>
      <Display lang="en"><![CDATA[Dell Express Flash NVMe PCIe SSD CM7 U.2 ISE Firmware Release]]></Display>
    </Name>
    <ComponentType value="FRMW">
      <Display lang="en"><![CDATA[Firmware]]></Display>
    </ComponentType>
    <Description>
      <Display lang="en"><![CDATA[Dell Express Flash NVMe PCIe SSD CM7 U.2 ISE Firmware Release]]></Display>
    </Description>
    <LUCategory value="Express Flash PCIe SSD">
      <Display lang="en"><![CDATA[Express Flash PCIe SSD]]></Display>
    </LUCategory>
    <Category value="Express Flash PCIe SSD">
      <Display lang="en"><![CDATA[Express Flash PCIe SSD]]></Display>
    </Category>
    <ImportantInfo URL="https://www.dell.com/support/home/en-us/drivers/DriversDetails?driverId=5JJFF">
      <Display lang="en"><![CDATA[NA]]></Display>
    </ImportantInfo>
    <SupportedDevices>
      <Device componentID="112167" embedded="0">
        <PCIInfo deviceID="0025" subDeviceID="223E" subVendorID="1028" vendorID="1E0F" />
        <Display lang="en"><![CDATA[Dell EMC PowerEdge Express Flash Ent NVMe AGN MU U.2 Gen4 1.6TB]]></Display>
        <RollbackInformation alternateRollbackIdentifier="112162" fmpWrapperIdentifier="55B4A079-CF96-4CB0-AC6F-DEB2238D9648" fmpWrapperVersion="1.0" impactsTPMmeasurements="true" rollbackIdentifier="A0BE5595-0612-450D-A1C7-159A25A0619D" rollbackTimeout="600" rollbackVolume="MAS022" />
        <PayloadConfiguration>
          <Image filename="*.bin" id="5192CAA7-85F4-4493-A5FD-F7EFDB9832BE" skip="false" type="FRMW" version="" />
        </PayloadConfiguration>
      </Device>
      <Device componentID="112163" embedded="0">
        <PCIInfo deviceID="0025" subDeviceID="2235" subVendorID="1028" vendorID="1E0F" />
        <Display lang="en"><![CDATA[Dell EMC PowerEdge Express Flash Ent NVMe AGN  RI U.2 Gen4 3.84TB]]></Display>
        <RollbackInformation alternateRollbackIdentifier="112162" fmpWrapperIdentifier="55B4A079-CF96-4CB0-AC6F-DEB2238D9648" fmpWrapperVersion="1.0" impactsTPMmeasurements="true" rollbackIdentifier="A0BE5595-0612-450D-A1C7-159A25A0619D" rollbackTimeout="600" rollbackVolume="MAS022" />
        <PayloadConfiguration>
          <Image filename="*.bin" id="5192CAA7-85F4-4493-A5FD-F7EFDB9832BE" skip="false" type="FRMW" version="" />
        </PayloadConfiguration>
      </Device>
      <Device componentID="112165" embedded="0">
        <PCIInfo deviceID="0025" subDeviceID="2233" subVendorID="1028" vendorID="1E0F" />
        <Display lang="en"><![CDATA[Dell EMC PowerEdge Express Flash Ent NVMe AGN  RI U.2 Gen4 15.36TB]]></Display>
        <RollbackInformation alternateRollbackIdentifier="112162" fmpWrapperIdentifier="55B4A079-CF96-4CB0-AC6F-DEB2238D9648" fmpWrapperVersion="1.0" impactsTPMmeasurements="true" rollbackIdentifier="A0BE5595-0612-450D-A1C7-159A25A0619D" rollbackTimeout="600" rollbackVolume="MAS022" />
        <PayloadConfiguration>
          <Image filename="*.bin" id="5192CAA7-85F4-4493-A5FD-F7EFDB9832BE" skip="false" type="FRMW" version="" />
        </PayloadConfiguration>
      </Device>
      <Device componentID="112164" embedded="0">
        <PCIInfo deviceID="0025" subDeviceID="2234" subVendorID="1028" vendorID="1E0F" />
        <Display lang="en"><![CDATA[Dell EMC PowerEdge Express Flash Ent NVMe AGN  RI U.2 Gen4 7.68TB]]></Display>
        <RollbackInformation alternateRollbackIdentifier="112162" fmpWrapperIdentifier="55B4A079-CF96-4CB0-AC6F-DEB2238D9648" fmpWrapperVersion="1.0" impactsTPMmeasurements="true" rollbackIdentifier="A0BE5595-0612-450D-A1C7-159A25A0619D" rollbackTimeout="600" rollbackVolume="MAS022" />
        <PayloadConfiguration>
          <Image filename="*.bin" id="5192CAA7-85F4-4493-A5FD-F7EFDB9832BE" skip="false" type="FRMW" version="" />
        </PayloadConfiguration>
      </Device>
      <Device componentID="112166" embedded="0">
        <PCIInfo deviceID="0025" subDeviceID="2232" subVendorID="1028" vendorID="1E0F" />
        <Display lang="en"><![CDATA[Dell EMC PowerEdge Express Flash Ent NVMe AGN  RI U.2 Gen4 30.72TB]]></Display>
        <RollbackInformation alternateRollbackIdentifier="112162" fmpWrapperIdentifier="55B4A079-CF96-4CB0-AC6F-DEB2238D9648" fmpWrapperVersion="1.0" impactsTPMmeasurements="true" rollbackIdentifier="A0BE5595-0612-450D-A1C7-159A25A0619D" rollbackTimeout="600" rollbackVolume="MAS022" />
        <PayloadConfiguration>
          <Image filename="*.bin" id="5192CAA7-85F4-4493-A5FD-F7EFDB9832BE" skip="false" type="FRMW" version="" />
        </PayloadConfiguration>
      </Device>
      <Device componentID="112162" embedded="0">
        <PCIInfo deviceID="0025" subDeviceID="2236" subVendorID="1028" vendorID="1E0F" />
        <Display lang="en"><![CDATA[Dell EMC PowerEdge Express Flash Ent NVMe AGN  RI U.2 Gen4 1.92TB]]></Display>
        <RollbackInformation alternateRollbackIdentifier="112162" fmpWrapperIdentifier="55B4A079-CF96-4CB0-AC6F-DEB2238D9648" fmpWrapperVersion="1.0" impactsTPMmeasurements="true" rollbackIdentifier="A0BE5595-0612-450D-A1C7-159A25A0619D" rollbackTimeout="600" rollbackVolume="MAS022" />
        <PayloadConfiguration>
          <Image filename="*.bin" id="5192CAA7-85F4-4493-A5FD-F7EFDB9832BE" skip="false" type="FRMW" version="" />
        </PayloadConfiguration>
      </Device>
      <Device componentID="112169" embedded="0">
        <PCIInfo deviceID="0025" subDeviceID="223C" subVendorID="1028" vendorID="1E0F" />
        <Display lang="en"><![CDATA[Dell EMC PowerEdge Express Flash Ent NVMe AGN  MU U.2 Gen4 6.4TB]]></Display>
        <RollbackInformation alternateRollbackIdentifier="112162" fmpWrapperIdentifier="55B4A079-CF96-4CB0-AC6F-DEB2238D9648" fmpWrapperVersion="1.0" impactsTPMmeasurements="true" rollbackIdentifier="A0BE5595-0612-450D-A1C7-159A25A0619D" rollbackTimeout="600" rollbackVolume="MAS022" />
        <PayloadConfiguration>
          <Image filename="*.bin" id="5192CAA7-85F4-4493-A5FD-F7EFDB9832BE" skip="false" type="FRMW" version="" />
        </PayloadConfiguration>
      </Device>
      <Device componentID="112168" embedded="0">
        <PCIInfo deviceID="0025" subDeviceID="223D" subVendorID="1028" vendorID="1E0F" />
        <Display lang="en"><![CDATA[Dell EMC PowerEdge Express Flash Ent NVMe AGN MU U.2 GEN4 3.2TB]]></Display>
        <RollbackInformation alternateRollbackIdentifier="112162" fmpWrapperIdentifier="55B4A079-CF96-4CB0-AC6F-DEB2238D9648" fmpWrapperVersion="1.0" impactsTPMmeasurements="true" rollbackIdentifier="A0BE5595-0612-450D-A1C7-159A25A0619D" rollbackTimeout="600" rollbackVolume="MAS022" />
        <PayloadConfiguration>
          <Image filename="*.bin" id="5192CAA7-85F4-4493-A5FD-F7EFDB9832BE" skip="false" type="FRMW" version="" />
        </PayloadConfiguration>
      </Device>
      <Device componentID="112170" embedded="0">
        <PCIInfo deviceID="0025" subDeviceID="223B" subVendorID="1028" vendorID="1E0F" />
        <Display lang="en"><![CDATA[Dell EMC PowerEdge Express Flash Ent NVMe AGN  MU U.2 Gen4 12.8TB]]></Display>
        <RollbackInformation alternateRollbackIdentifier="112162" fmpWrapperIdentifier="55B4A079-CF96-4CB0-AC6F-DEB2238D9648" fmpWrapperVersion="1.0" impactsTPMmeasurements="true" rollbackIdentifier="A0BE5595-0612-450D-A1C7-159A25A0619D" rollbackTimeout="600" rollbackVolume="MAS022" />
        <PayloadConfiguration>
          <Image filename="*.bin" id="5192CAA7-85F4-4493-A5FD-F7EFDB9832BE" skip="false" type="FRMW" version="" />
        </PayloadConfiguration>
      </Device>
    </SupportedDevices>
    <RevisionHistory>
      <Display lang="en"><![CDATA[Dell Express Flash NVMe PCIe SSD CM7 U.2 ISE Firmware Release]]></Display>
    </RevisionHistory>
    <Criticality value="1">
      <Display lang="en"><![CDATA[Recommended]]></Display>
    </Criticality>
    <SupportedSystems>
      <Brand key="3" prefix="PE">
        <Display lang="en"><![CDATA[PowerEdge]]></Display>
        <Model systemID="0B4F" systemIDType="BIOS">
          <Display lang="en"><![CDATA[HS5610]]></Display>
        </Model>
        <Model systemID="08FC" systemIDType="BIOS">
          <Display lang="en"><![CDATA[R6515]]></Display>
        </Model>
        <Model systemID="0B9D" systemIDType="BIOS">
          <Display lang="en"><![CDATA[R660XS]]></Display>
        </Model>
      </Brand>
    </SupportedSystems>
    <FMPWrappers>
      <FMPWrapperInformation digitalSignature="false" driverFileName="DellNVMeSSDFWUpdDrv" filePathName="FmpUpdateWrapper.efi" identifier="55B4A079-CF96-4CB0-AC6F-DEB2238D9648" name="PCIeSSD" vendorCodeType="UINT32">
        <Inventory source="Device" supported="true" />
        <Update rollback="true" supported="true" />
      </FMPWrapperInformation>
    </FMPWrappers>
  </SoftwareComponent>
</Manifest>`

// setupMockCatalogFile creates a temporary file with mock catalog data and returns the file path
// The caller is responsible for removing the file when done
func setupMockCatalogFile() (string, error) {
	tmpFile, err := os.CreateTemp("", "mock-catalog-*.xml")
	if err != nil {
		return "", err
	}

	_, err = tmpFile.WriteString(mockCatalogXML)
	if err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return "", err
	}

	tmpFile.Close()
	return tmpFile.Name(), nil
}

func TestProcessDellCatalog_Local(t *testing.T) {
	ctx := context.Background()

	// Create a temporary mock catalog file
	mockCatalogPath, err := setupMockCatalogFile()
	if err != nil {
		t.Fatalf("Failed to create mock catalog file: %v", err)
	}
	defer os.Remove(mockCatalogPath) // Clean up after test

	vendorCatalog := &VendorCatalog{
		VendorLocalCatalogPath: mockCatalogPath,
	}

	err = vendorCatalog.processDellCatalog(ctx)

	if err != nil {
		t.Errorf("processDellCatalog() returned an error: %v", err)
	}

	if vendorCatalog.CatalogInfo.VendorId == nil || *vendorCatalog.CatalogInfo.VendorId == "" {
		t.Errorf("processDellCatalog() did not set VendorId")
	}
	if vendorCatalog.CatalogInfo.VendorReleaseTimestamp == nil || *vendorCatalog.CatalogInfo.VendorReleaseTimestamp == "" {
		t.Errorf("processDellCatalog() did not set VendorReleaseTimestamp")
	}
	if len(vendorCatalog.Binaries) == 0 {
		t.Errorf("processDellCatalog() did not populate Binaries")
	}
}

func TestProcessDellCatalog_Filtered(t *testing.T) {
	ctx := context.Background()

	// Create a temporary mock catalog file
	mockCatalogPath, err := setupMockCatalogFile()
	if err != nil {
		t.Fatalf("Failed to create mock catalog file: %v", err)
	}
	defer os.Remove(mockCatalogPath) // Clean up after test

	vendorCatalog := &VendorCatalog{
		VendorLocalCatalogPath: mockCatalogPath,
		VendorSystemsFilter:    []string{"PowerEdge R730"},
	}

	err = vendorCatalog.processDellCatalog(ctx)

	if err != nil {
		t.Errorf("processDellCatalog() returned an error: %v", err)
	}

	if vendorCatalog.CatalogInfo.VendorId == nil || *vendorCatalog.CatalogInfo.VendorId == "" {
		t.Errorf("processDellCatalog() did not set VendorId")
	}
	if vendorCatalog.CatalogInfo.VendorReleaseTimestamp == nil || *vendorCatalog.CatalogInfo.VendorReleaseTimestamp == "" {
		t.Errorf("processDellCatalog() did not set VendorReleaseTimestamp")
	}
	if len(vendorCatalog.Binaries) == 0 {
		t.Errorf("processDellCatalog() did not populate Binaries")
	}
}

func TestProcessDellCatalog_FilteredByModelOnly(t *testing.T) {
	ctx := context.Background()

	// Create a temporary mock catalog file
	mockCatalogPath, err := setupMockCatalogFile()
	if err != nil {
		t.Fatalf("Failed to create mock catalog file: %v", err)
	}
	defer os.Remove(mockCatalogPath)

	// Filter using just the model name "R430" instead of "PowerEdge R430"
	vendorCatalog := &VendorCatalog{
		VendorLocalCatalogPath: mockCatalogPath,
		VendorSystemsFilter:    []string{"R430"},
	}

	err = vendorCatalog.processDellCatalog(ctx)

	if err != nil {
		t.Errorf("processDellCatalog() returned an error: %v", err)
	}

	if len(vendorCatalog.Binaries) == 0 {
		t.Errorf("processDellCatalog() should find binaries when filtering by model name only")
	}

	// Verify the correct binary was matched (the first SoftwareComponent which supports R430)
	if len(vendorCatalog.Binaries) != 1 {
		t.Errorf("processDellCatalog() expected 1 binary, got %d", len(vendorCatalog.Binaries))
	}
}

func TestProcessDellCatalog_WrongLocalPath(t *testing.T) {
	ctx := context.Background()

	vendorCatalog := &VendorCatalog{
		VendorLocalCatalogPath: "./missing_path/Catalog.xml",
	}

	err := vendorCatalog.processDellCatalog(ctx)

	if err == nil {
		t.Errorf("processDellCatalog() should have returned an error for missing path")
	}

	if !strings.Contains(err.Error(), "no such file or directory") {
		t.Errorf("processDellCatalog() returned unexpected error: %v", err)
	}
}
