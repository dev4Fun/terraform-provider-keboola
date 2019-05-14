## 0.1.3 (Unreleased)

IMPROVEMENTS:

* `provider`: Whitespace is now trimmed from the API key when it is read from configuration, which will help alleviate issues on first set up of terraform.

FIXES:

* `keboola_snowflake_writer`: `ID` and `workspaceID` on `ProvisionSnowflakeResponse` are now `string` types instead of `int` types.
* `keboola_snowflake_writer_tables`: Data filtering columns (e.g. `where_column`, `changed_since`) are now mapped.

## 0.1.4 (Unreleased)

* `Keboola_storage_bucket`: Link shared buckets between projects. Using `is_linked`, `source_project_id`, `source_bucket_id`
