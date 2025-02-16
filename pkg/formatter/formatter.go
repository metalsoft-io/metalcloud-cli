package formatter

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"reflect"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

const (
	ConfigFormat = "format"
)

type RecordFieldConfig struct {
	Title    string
	Hidden   bool
	IsStatus bool
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

	var columnConfigs []table.ColumnConfig

	// Check if the result is a struct with slice field named data
	paginatedData, ok := getPaginatedData(result)
	if ok {
		// Loop through the paginated data list
		headerSet := false
		for i := 0; i < paginatedData.Len(); i++ {
			names, values, configs := getFieldNamesAndValues(paginatedData.Index(i).Interface(), printConfig)
			if !headerSet {
				t.AppendHeader(names)
				columnConfigs = configs

				headerSet = true
			}
			t.AppendRows([]table.Row{values})
		}
	} else {
		// Print the result as a table
		names, values, configs := getFieldNamesAndValues(result, printConfig)
		t.AppendHeader(names)
		t.AppendRows([]table.Row{values})
		columnConfigs = configs
	}

	t.SetStyle(table.StyleLight)
	if len(columnConfigs) > 0 {
		t.SetColumnConfigs(columnConfigs)
	}

	return t
}

func getPaginatedData(record interface{}) (*reflect.Value, bool) {
	// Get the reflected value and type
	val := reflect.ValueOf(record)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Check if it's a struct
	if val.Kind() != reflect.Struct {
		return nil, false
	}

	// Look for "Data" field
	field := val.FieldByName("Data")
	if !field.IsValid() {
		return nil, false
	}

	// Check if field is a slice
	if field.Kind() != reflect.Slice {
		return nil, false
	}

	return &field, true
}

func getFieldNamesAndValues(record interface{}, printConfig *PrintConfig) (table.Row, table.Row, []table.ColumnConfig) {
	// Handle nil input
	if record == nil {
		return nil, nil, nil
	}

	// Get reflected value
	val := reflect.ValueOf(record)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Must be a struct
	if val.Kind() != reflect.Struct {
		return nil, nil, nil
	}

	// Get type for field names
	fieldType := val.Type()
	fieldCount := val.NumField()

	names := make(table.Row, 0)
	values := make(table.Row, 0)
	configs := make([]table.ColumnConfig, 0)

	// Iterate through fields
	for i := 0; i < fieldCount; i++ {
		field := fieldType.Field(i)

		// Check if the print config exists
		if printConfig != nil {
			fieldConfig, ok := printConfig.FieldsConfig[field.Name]
			if ok && !fieldConfig.Hidden {
				title := field.Name
				if fieldConfig.Title != "" {
					title = fieldConfig.Title
				}
				names = append(names, title)
				values = append(values, val.Field(i).Interface())

				if fieldConfig.IsStatus {
					configs = append(configs, table.ColumnConfig{
						Name:        fieldConfig.Title,
						Transformer: text.Transformer(formatStatusValue),
					})
				}
			}
		} else {
			names = append(names, field.Name)
			values = append(values, val.Field(i).Interface())
		}
	}

	return names, values, configs
}

func formatStatusValue(value interface{}) string {
	if _, ok := value.(string); ok {
		var color text.Color
		switch value {
		case "available":
			color = text.FgBlue
		case "used":
			color = text.FgGreen
		case "unavailable":
			color = text.FgMagenta
		default:
			color = text.FgYellow
		}
		return color.Sprintf("%s", value)
	}
	return fmt.Sprint(value)
}
