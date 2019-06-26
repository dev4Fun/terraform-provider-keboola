package keboola

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

type FTPFile struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Configuration string `json:"configuration`
	IsDisabled    bool   `json:"isDisabled"`
}

func resourceKeboolaFTPExtractorFiles() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeboolaFTPExtractorFilesCreate,
		// Read:   resourceKeboolaFTPExtractorFilesRead,
		// Update: resourceKeboolaFTPExtractorFilesUpdate,
		// Delete: resourceKeboolaFTPExtractorFilesDelete,

		Schema: map[string]*schema.Schema{
			"extractor_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"files": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"configuration": {
							Type:     schema.TypeString,
							Required: true,
						},
						"is_disabled": {
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceKeboolaFTPExtractorFilesCreate(d *schema.ResourceData, meta interface{}) error {
	log.Println("[INFO] Creating FTP Extractor Files in Keboola.")

	extractorID := d.Get("extractor_id").(string)
	files := d.Get("files").([]interface{})

	mappedFiles := make([]FTPFile, 0, len(files))

	for _, file := range files {
		config := file.(map[string]interface{})

		mappedFile := FTPFile{
			Name:          config["name"].(string),
			Description:   config["description"].(string),
			Configuration: config["configuration"].(string),
			IsDisabled:    config["is_disabled"].(bool),
		}

		mappedFiles = append(mappedFiles, mappedFile)
	}

	client := meta.(*KBCClient)

	getExtractorResponse, err := client.GetFromStorage(fmt.Sprintf("storage/components/keboola.ex-ftp/configs/%s", extractorID))

	if hasErrors(err, getExtractorResponse) {
		return extractError(err, getExtractorResponse)
	}

	var ftpExtractor FTPExtractor

	decoder := json.NewDecoder(getExtractorResponse.Body)
	err = decoder.Decode(&ftpExtractor)

	if err != nil {
		return err
	}

	ftpExtractor.Files = mappedFiles

	ftpConfigJSON, err := json.Marshal(ftpExtractor)

	if err != nil {
		return err
	}

	// updatePostgreSQLForm := url.Values{}
	// updatePostgreSQLForm.Add("configuration", string(postgresqlConfigJSON))
	// updatePostgreSQLForm.Add("change_description", "Update PostgreSQL tables")

	// updatePostgreSQLBuffer := buffer.FromForm(updatePostgreSQLForm)

	// updateResponse, err := client.PutToStorage(fmt.Sprintf("storage/components/keboola.wr-db-pgsql/configs/%s", writerID), updatePostgreSQLBuffer)

	// if hasErrors(err, updateResponse) {
	// 	return extractError(err, updateResponse)
	// }

	// d.SetId(writerID)

	// return resourceKeboolaFTPExtractorFilesRead(d, meta)
}

// func resourceKeboolaFTPExtractorFilesRead(d *schema.ResourceData, meta interface{}) error {
// 	log.Println("[INFO] Reading PostgreSQL Writer Tables from Keboola.")

// 	if d.Id() == "" {
// 		return nil
// 	}

// 	client := meta.(*KBCClient)

// 	getPostgreSQLWriterResponse, err := client.GetFromStorage(fmt.Sprintf("storage/components/keboola.wr-db-pgsql/configs/%s", d.Id()))

// 	if hasErrors(err, getPostgreSQLWriterResponse) {
// 		if getPostgreSQLWriterResponse.StatusCode == 404 {
// 			d.SetId("")
// 			return nil
// 		}

// 		return extractError(err, getPostgreSQLWriterResponse)
// 	}

// 	var postgresqlWriter PostgreSQLWriter

// 	decoder := json.NewDecoder(getPostgreSQLWriterResponse.Body)
// 	err = decoder.Decode(&postgresqlWriter)

// 	if err != nil {
// 		return err
// 	}

// 	var tables []map[string]interface{}

// 	for _, tableConfig := range postgresqlWriter.Configuration.Parameters.Tables {
// 		tableDetails := map[string]interface{}{
// 			"db_name":     tableConfig.DatabaseName,
// 			"export":      tableConfig.Export,
// 			"table_id":    tableConfig.TableID,
// 			"incremental": tableConfig.Incremental,
// 			"primary_key": tableConfig.PrimaryKey,
// 		}

// 		var columns []map[string]interface{}

// 		for _, item := range tableConfig.Items {
// 			columnDetails := map[string]interface{}{
// 				"name":     item.Name,
// 				"db_name":  item.DatabaseName,
// 				"type":     item.Type,
// 				"size":     item.Size,
// 				"nullable": item.IsNullable,
// 				"default":  item.DefaultValue,
// 			}

// 			columns = append(columns, columnDetails)
// 		}

// 		tableDetails["column"] = columns

// 		tables = append(tables, tableDetails)
// 	}

// 	d.Set("table", tables)

// 	return nil
// }

// func resourceKeboolaFTPExtractorFilesUpdate(d *schema.ResourceData, meta interface{}) error {
// 	log.Println("[INFO] Updating PostgreSQL Writer Tables in Keboola.")

// 	tables := d.Get("table").([]interface{})

// 	mappedTables := make([]PostgreSQLWriterTable, 0, len(tables))
// 	storageTables := make([]PostgreSQLWriterStorageTable, 0, len(tables))

// 	for _, table := range tables {
// 		config := table.(map[string]interface{})

// 		mappedTable := PostgreSQLWriterTable{
// 			DatabaseName: config["db_name"].(string),
// 			Export:       config["export"].(bool),
// 			TableID:      config["table_id"].(string),
// 			Incremental:  config["incremental"].(bool),
// 		}

// 		if q := config["primary_key"]; q != nil {
// 			mappedTable.PrimaryKey = AsStringArray(q.([]interface{}))
// 		}

// 		storageTable := PostgreSQLWriterStorageTable{
// 			Source:      mappedTable.TableID,
// 			Destination: fmt.Sprintf("%s.csv", mappedTable.TableID),
// 		}

// 		columnConfigs := config["column"].([]interface{})
// 		mappedColumns := make([]PostgreSQLWriterTableItem, 0, len(columnConfigs))
// 		columnNames := make([]string, 0, len(columnConfigs))
// 		for _, column := range columnConfigs {
// 			columnConfig := column.(map[string]interface{})

// 			mappedColumn := PostgreSQLWriterTableItem{
// 				Name:         columnConfig["name"].(string),
// 				DatabaseName: columnConfig["db_name"].(string),
// 				Type:         columnConfig["type"].(string),
// 				Size:         columnConfig["size"].(string),
// 				IsNullable:   columnConfig["nullable"].(bool),
// 				DefaultValue: columnConfig["default"].(string),
// 			}

// 			mappedColumns = append(mappedColumns, mappedColumn)
// 			columnNames = append(columnNames, mappedColumn.Name)
// 		}

// 		mappedTable.Items = mappedColumns
// 		storageTable.Columns = columnNames

// 		mappedTables = append(mappedTables, mappedTable)
// 		storageTables = append(storageTables, storageTable)
// 	}

// 	client := meta.(*KBCClient)

// 	getWriterResponse, err := client.GetFromStorage(fmt.Sprintf("storage/components/keboola.wr-db-pgsql/configs/%s", d.Id()))

// 	if hasErrors(err, getWriterResponse) {
// 		return extractError(err, getWriterResponse)
// 	}

// 	var postgresqlWriter PostgreSQLWriter

// 	decoder := json.NewDecoder(getWriterResponse.Body)
// 	err = decoder.Decode(&postgresqlWriter)

// 	if err != nil {
// 		return err
// 	}

// 	postgresqlWriter.Configuration.Parameters.Tables = mappedTables
// 	postgresqlWriter.Configuration.Storage.Input.Tables = storageTables

// 	postgresqlConfigJSON, err := json.Marshal(postgresqlWriter.Configuration)

// 	if err != nil {
// 		return err
// 	}

// 	updatePostgreSQLForm := url.Values{}
// 	updatePostgreSQLForm.Add("configuration", string(postgresqlConfigJSON))
// 	updatePostgreSQLForm.Add("changeDescription", "Update PostgreSQL tables")

// 	updatePostgreSQLBuffer := buffer.FromForm(updatePostgreSQLForm)

// 	updateResponse, err := client.PutToStorage(fmt.Sprintf("storage/components/keboola.wr-db-pgsql/configs/%s", d.Id()), updatePostgreSQLBuffer)

// 	if hasErrors(err, updateResponse) {
// 		return extractError(err, updateResponse)
// 	}

// 	return resourceKeboolaFTPExtractorFilesRead(d, meta)
// }

// func resourceKeboolaFTPExtractorFilesDelete(d *schema.ResourceData, meta interface{}) error {
// 	log.Printf("[INFO] Clearing PostgreSQL Writer Tables in Keboola: %s", d.Id())

// 	client := meta.(*KBCClient)

// 	getWriterResponse, err := client.GetFromStorage(fmt.Sprintf("storage/components/keboola.wr-db-pgsql/configs/%s", d.Id()))

// 	if hasErrors(err, getWriterResponse) {
// 		return extractError(err, getWriterResponse)
// 	}

// 	var postgresqlWriter PostgreSQLWriter

// 	decoder := json.NewDecoder(getWriterResponse.Body)
// 	err = decoder.Decode(&postgresqlWriter)

// 	if err != nil {
// 		return err
// 	}

// 	var emptyTables []PostgreSQLWriterTable
// 	postgresqlWriter.Configuration.Parameters.Tables = emptyTables

// 	var emptyStorageTables []PostgreSQLWriterStorageTable
// 	postgresqlWriter.Configuration.Storage.Input.Tables = emptyStorageTables

// 	postgresqlConfigJSON, err := json.Marshal(postgresqlWriter.Configuration)

// 	if err != nil {
// 		return err
// 	}

// 	clearPostgreSQLTablesForm := url.Values{}
// 	clearPostgreSQLTablesForm.Add("configuration", string(postgresqlConfigJSON))
// 	clearPostgreSQLTablesForm.Add("changeDescription", "Update PostgreSQL tables")

// 	clearPostgreSQLTablesBuffer := buffer.FromForm(clearPostgreSQLTablesForm)

// 	clearResponse, err := client.PutToStorage(fmt.Sprintf("storage/components/keboola.wr-db-pgsql/configs/%s", d.Id()), clearPostgreSQLTablesBuffer)

// 	if hasErrors(err, clearResponse) {
// 		return extractError(err, clearResponse)
// 	}

// 	d.SetId("")

// 	return nil
// }
