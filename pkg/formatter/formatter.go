package formatter

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"reflect"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	sdk "github.com/metalsoft-io/metalcloud-sdk-go"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

const (
	ConfigFormat = "format"
)

type RecordFieldConfig struct {
	Title       string
	Hidden      bool
	Order       int
	MaxWidth    int
	Transformer func(interface{}) string
	InnerFields map[string]RecordFieldConfig
}

type PrintConfig struct {
	FieldsConfig map[string]RecordFieldConfig
}

func PrintResult(result interface{}, printConfig *PrintConfig) error {
	format := strings.ToLower(viper.GetString(ConfigFormat))

	switch format {
	case "json":
		// Print JSON result
		jsonResult, err := json.Marshal(result)
		if err != nil {
			return fmt.Errorf("failed to convert to JSON: %v", err)
		}
		fmt.Printf("%s", string(jsonResult))
	case "yaml":
		// Convert JSON to YAML
		yamlResult, err := yaml.Marshal(result)
		if err != nil {
			return fmt.Errorf("failed to convert to YAML: %v", err)
		}
		fmt.Printf("%s", string(yamlResult))
	case "text":
		generateTable(result, printConfig).Render()
	case "csv":
		generateTable(result, printConfig).RenderCSV()
	case "md":
		generateTable(result, printConfig).RenderMarkdown()
	default:
		return fmt.Errorf("%s format not supported yet", format)
	}

	return nil
}

func generateTable(result interface{}, printConfig *PrintConfig) table.Writer {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	// Check if the result is a struct with slice field named data
	paginatedData, ok := getPaginatedData(result)
	if ok {
		// Loop through the paginated data list
		headerSet := false
		for i := 0; i < paginatedData.Len(); i++ {
			names, values, configs := getFieldNamesAndValues(paginatedData.Index(i).Interface(), printConfig)
			if !headerSet {
				t.AppendHeader(names)
				t.SetColumnConfigs(configs)

				headerSet = true
			}
			t.AppendRows([]table.Row{values})
		}
	} else if reflect.TypeOf(result).Kind() == reflect.String {
		// Check if the result is a string
		t.AppendHeader(table.Row{"Result"})
		t.AppendRows([]table.Row{{result}})
	} else if reflect.TypeOf(result).Kind() == reflect.Map {
		// Check if the result is a map
		keys := make(table.Row, 0)
		values := make(table.Row, 0)
		for k, v := range result.(map[string]interface{}) {
			keys = append(keys, k)
			values = append(values, v)
		}
		t.AppendHeader(keys)
		t.AppendRows([]table.Row{values})
	} else {
		// Print the result as a table
		names, values, configs := getFieldNamesAndValues(result, printConfig)
		t.AppendHeader(names)
		t.AppendRows([]table.Row{values})
		t.SetColumnConfigs(configs)
	}

	t.SetStyle(table.StyleLight)

	return t
}

// getPaginatedData checks if the record is a slice or struct with a slice field named Data
// and returns the value of the slice if it exists
func getPaginatedData(record interface{}) (*reflect.Value, bool) {
	recordValue := reflect.ValueOf(record)
	if recordValue.Kind() == reflect.Ptr {
		recordValue = recordValue.Elem()
	}

	if recordValue.Kind() == reflect.Slice {
		return &recordValue, true
	}

	if recordValue.Kind() != reflect.Struct {
		return nil, false
	}

	dataField := recordValue.FieldByName("Data")
	if !dataField.IsValid() {
		return nil, false
	}

	if dataField.Kind() != reflect.Slice {
		return nil, false
	}

	return &dataField, true
}

func getFieldNamesAndValues(record interface{}, printConfig *PrintConfig) (table.Row, table.Row, []table.ColumnConfig) {
	if record == nil {
		return nil, nil, nil
	}

	recordValue := reflect.ValueOf(record)
	if recordValue.Kind() == reflect.Ptr {
		recordValue = recordValue.Elem()
	}

	if recordValue.Kind() != reflect.Struct {
		return nil, nil, nil
	}

	fieldType := recordValue.Type()

	names := make(table.Row, 0)
	values := make(table.Row, 0)
	configs := make([]table.ColumnConfig, 0)

	if printConfig == nil {
		for i := 0; i < recordValue.NumField(); i++ {
			names = append(names, fieldType.Field(i).Name)
			values = append(values, extractValue(recordValue.Field(i)))
		}
	} else {
		count, maxOrder := getColumnsCount(&printConfig.FieldsConfig)
		if count < maxOrder {
			count = maxOrder
		}

		names = make(table.Row, count)
		values = make(table.Row, count)

		populate(record, &printConfig.FieldsConfig, &names, &values, &configs)
	}

	return names, values, configs
}

func getColumnsCount(fieldsConfig *map[string]RecordFieldConfig) (int, int) {
	count := 0
	maxOrder := 0

	for fieldName, fieldConfig := range *fieldsConfig {
		if !fieldConfig.Hidden {
			count++
			if fieldConfig.Order > maxOrder {
				maxOrder = fieldConfig.Order
			}
		}

		if len(fieldConfig.InnerFields) > 0 {
			innerCount, innerMaxOrder := getColumnsCount(&fieldConfig.InnerFields)

			count += innerCount
			if innerMaxOrder > maxOrder {
				maxOrder = innerMaxOrder
			}
		}

		if fieldConfig.Order == 0 {
			tempFieldConfig := (*fieldsConfig)[fieldName]
			tempFieldConfig.Order = count
			(*fieldsConfig)[fieldName] = tempFieldConfig
		}
	}

	return count, maxOrder
}

func populate(record interface{}, fieldsConfig *map[string]RecordFieldConfig, names *table.Row, values *table.Row, configs *[]table.ColumnConfig) {
	if record == nil {
		return
	}

	recordValue := reflect.ValueOf(record)
	if recordValue.Kind() == reflect.Ptr {
		recordValue = recordValue.Elem()
	}

	if recordValue.Kind() != reflect.Struct {
		// Add headers even if the record is not a valid value
		for fieldName, fieldConfig := range *fieldsConfig {
			addField(fieldConfig, fieldName, nil, names, values, configs)
		}
		return
	}

	fieldType := recordValue.Type()

	for i := 0; i < recordValue.NumField(); i++ {
		field := fieldType.Field(i)

		fieldConfig, ok := (*fieldsConfig)[field.Name]
		if ok {
			addField(fieldConfig, field.Name, extractValue(recordValue.Field(i)), names, values, configs)

			if len(fieldConfig.InnerFields) > 0 {
				populate(recordValue.Field(i).Interface(), &fieldConfig.InnerFields, names, values, configs)
			}
		}
	}
}

func addField(fieldConfig RecordFieldConfig, fieldName string, fieldValue interface{}, names *table.Row, values *table.Row, configs *[]table.ColumnConfig) {
	if !fieldConfig.Hidden {
		title := fieldName
		if fieldConfig.Title != "" {
			title = fieldConfig.Title
		}

		(*names)[fieldConfig.Order-1] = title
		if fieldConfig.Transformer != nil || fieldConfig.MaxWidth > 0 {
			*configs = append(*configs, table.ColumnConfig{
				Name:        title,
				WidthMax:    fieldConfig.MaxWidth,
				Transformer: text.Transformer(fieldConfig.Transformer),
			})
		}

		(*values)[fieldConfig.Order-1] = fieldValue
	}
}

func extractValue(value reflect.Value) interface{} {
	if !value.IsValid() {
		return nil
	}

	if value.Kind() == reflect.Pointer {
		value = value.Elem()
	}

	switch value.Kind() {
	case reflect.String:
		return value.String()
	case reflect.Bool:
		return value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return value.Uint()
	case reflect.Float32, reflect.Float64:
		return value.Float()
	case reflect.Array, reflect.Slice:
		var result []interface{}
		for i := 0; i < value.Len(); i++ {
			result = append(result, extractValue(value.Index(i)))
		}
		return result
	case reflect.Map:
		var result []string
		for _, key := range value.MapKeys() {
			result = append(result, fmt.Sprintf("%s: %s", key.String(), extractValue(value.MapIndex(key))))
		}
		return result
	case reflect.Struct:
		if reflect.TypeOf(value) == reflect.TypeOf(sdk.NullableInt32{}) {
			if value.Interface().(sdk.NullableInt32).IsSet() {
				return value.Interface().(sdk.NullableInt32).Get()
			}
		}
		return nil
	case reflect.Invalid:
		return nil
	default:
		return value.String()
	}
}

func FormatStatusValue(value interface{}) string {
	if _, ok := value.(string); ok {
		var color text.Color
		switch value {
		case "available":
			color = text.FgBlue
		case "ready":
			color = text.FgBlue
		case "used":
			color = text.FgGreen
		case "unavailable":
			color = text.FgMagenta
		case "registering":
			color = text.FgCyan
		case "cleaning":
			color = text.FgHiCyan
		case "cleaning_required":
			color = text.FgHiCyan
		case "updating_firmware":
			color = text.FgHiYellow
		case "pending_registration":
			color = text.FgHiYellow
		case "used_registering":
			color = text.FgHiGreen
		case "used_diagnostics":
			color = text.FgHiGreen
		case "decommissioned":
			color = text.FgHiRed
		case "removed_from_rack":
			color = text.FgHiRed
		case "defective":
			color = text.FgRed
		case "active":
			color = text.FgGreen
		case "ordered":
			color = text.FgCyan
		case "draft":
			color = text.FgCyan
		default:
			color = text.FgYellow
		}
		return color.Sprintf("%s", value)
	}
	return fmt.Sprint(value)
}

func FormatDateTimeValue(value interface{}) string {
	if value == nil {
		return ""
	}

	if stringValue, ok := value.(string); ok {
		tm, err := time.Parse("2006-01-02T15:04:05Z", stringValue)
		if err == nil {
			return tm.Local().Format(time.RFC822)
		}

		tm, err = time.Parse("2006-01-02T15:04:05.000Z", stringValue)
		if err == nil {
			return tm.Local().Format(time.RFC822)
		}

		tm, err = time.Parse("2006-01-02T15:04:05.", stringValue)
		if err == nil {
			return tm.Local().Format(time.RFC822)
		}

		if stringValue == "0000-00-00T00:00:00Z" {
			return ""
		}
	}

	if tm, ok := value.(time.Time); ok {
		return tm.Local().Format(time.RFC822)
	}

	return fmt.Sprint(value)
}
