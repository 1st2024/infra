// Define an environment named "local"
env "local" {
  // Declare where the schema definition resides.
  // Also supported: ["file://multi.hcl", "file://schema.hcl"].
  src = "ent://pkg/schema/"

  // Define the URL of the Dev Database for this environment
  // See: https://atlasgo.io/concepts/dev-database
  dev = "docker://postgres/15/dev?search_path=public"

  migration {
    // Define the path to the migration directory.
    // See: https://entgo.io/docs/migrate/#migration-directory
    dir = "file://migrations"
  }

  diff  {
  	skip  {
  		drop_func = true
        drop_trigger = true
        drop_schema = true
  	}
  }
}
