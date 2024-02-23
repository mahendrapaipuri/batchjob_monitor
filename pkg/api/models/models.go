// Package models defines different models used in stats
package models

import "github.com/mahendrapaipuri/ceems/internal/structset"

const (
	unitsTableName = "units"
	usageTableName = "usage"
)

// Unit is an abstract compute unit that can mean Job (batchjobs), VM (cloud) or Pod (k8s)
type Unit struct {
	ID                  int64      `json:"-"                          sql:"id"                         sqlitetype:"integer not null primary key"`
	UUID                string     `json:"uuid"                       sql:"uuid"                       sqlitetype:"text"`              // Unique identifier of unit. It can be Job ID for batch jobs, UUID for pods in k8s or VMs in Openstack
	Name                string     `json:"name"                       sql:"name"                       sqlitetype:"text"`              // Name of compute unit
	Project             string     `json:"project"                    sql:"project"                    sqlitetype:"text"`              // Account in batch systems, Tenant in Openstack, Namespace in k8s
	Grp                 string     `json:"grp"                        sql:"grp"                        sqlitetype:"text"`              // User group
	Usr                 string     `json:"usr"                        sql:"usr"                        sqlitetype:"text"`              // Username
	Submit              string     `json:"submit"                     sql:"submit"                     sqlitetype:"text"`              // Submission time
	Start               string     `json:"start"                      sql:"start"                      sqlitetype:"text"`              // Start time
	End                 string     `json:"end"                        sql:"end"                        sqlitetype:"text"`              // End time
	SubmitTS            int64      `json:"submit_ts"                  sql:"submit_ts"                  sqlitetype:"integer"`           // Submission timestamp
	StartTS             int64      `json:"start_ts"                   sql:"start_ts"                   sqlitetype:"integer"`           // Start timestamp
	EndTS               int64      `json:"end_ts"                     sql:"end_ts"                     sqlitetype:"integer"`           // End timestamp
	Elapsed             string     `json:"elapsed"                    sql:"elapsed"                    sqlitetype:"text"`              // Total elapsed time
	ElapsedRaw          int64      `json:"elapsed_raw"                sql:"elapsed_raw"                sqlitetype:"integer"`           // Total elapsed time in seconds
	Exitcode            string     `json:"exitcode"                   sql:"exitcode"                   sqlitetype:"text"`              // Exit code of unit
	State               string     `json:"state"                      sql:"state"                      sqlitetype:"text"`              // Current state of unit
	Allocation          Allocation `json:"allocation"                 sql:"allocation"                 sqlitetype:"text default '{}'"` // Allocation map of unit. Only string and int64 values are supported in map
	TotalCPUBilling     int64      `json:"total_cpu_billing"          sql:"total_cpu_billing"          sqlitetype:"integer"`           // Total CPU billing for unit
	TotalGPUBilling     int64      `json:"total_gpu_billing"          sql:"total_gpu_billing"          sqlitetype:"integer"`           // Total GPU billing for unit
	TotalMiscBilling    int64      `json:"total_misc_billing"         sql:"total_misc_billing"         sqlitetype:"integer"`           // Total billing for unit that are not in CPU and GPU billing
	AveCPUUsage         float64    `json:"avg_cpu_usage"              sql:"avg_cpu_usage"              sqlitetype:"real"`              // Average CPU usage during lifetime of unit
	AveCPUMemUsage      float64    `json:"avg_cpu_mem_usage"          sql:"avg_cpu_mem_usage"          sqlitetype:"real"`              // Average CPU memory during lifetime of unit
	TotalCPUEnergyUsage float64    `json:"total_cpu_energy_usage_kwh" sql:"total_cpu_energy_usage_kwh" sqlitetype:"real"`              // Total CPU energy usage in kWh during lifetime of unit
	TotalCPUEmissions   float64    `json:"total_cpu_emissions_gms"    sql:"total_cpu_emissions_gms"    sqlitetype:"real"`              // Total CPU emissions in grams during lifetime of unit
	AveGPUUsage         float64    `json:"avg_gpu_usage"              sql:"avg_gpu_usage"              sqlitetype:"real"`              // Average GPU usage during lifetime of unit
	AveGPUMemUsage      float64    `json:"avg_gpu_mem_usage"          sql:"avg_gpu_mem_usage"          sqlitetype:"real"`              // Average GPU memory during lifetime of unit
	TotalGPUEnergyUsage float64    `json:"total_gpu_energy_usage_kwh" sql:"total_gpu_energy_usage_kwh" sqlitetype:"real"`              // Total GPU energy usage in kWh during lifetime of unit
	TotalGPUEmissions   float64    `json:"total_gpu_emissions_gms"    sql:"total_gpu_emissions_gms"    sqlitetype:"real"`              // Total GPU emissions in grams during lifetime of unit
	TotalIOWriteHot     float64    `json:"total_io_write_hot_gb"      sql:"total_io_write_hot_gb"      sqlitetype:"real"`              // Total IO write on hot storage in GB during lifetime of unit
	TotalIOReadHot      float64    `json:"total_io_read_hot_gb"       sql:"total_io_read_hot_gb"       sqlitetype:"real"`              // Total IO read on hot storage in GB during lifetime of unit
	TotalIOWriteCold    float64    `json:"total_io_write_cold_gb"     sql:"total_io_write_cold_gb"     sqlitetype:"real"`              // Total IO write on cold storage in GB during lifetime of unit
	TotalIOReadCold     float64    `json:"total_io_read_cold_gb"      sql:"total_io_read_cold_gb"      sqlitetype:"real"`              // Total IO read on cold storage in GB during lifetime of unit
	TotalIngress        float64    `json:"total_ingress_in_gb"        sql:"total_ingress_in_gb"        sqlitetype:"real"`              // Total ingress traffic in GB of unit
	TotalOutgress       float64    `json:"total_outgress_in_gb"       sql:"total_outgress_in_gb"       sqlitetype:"real"`              // Total outgress traffic in GB of unit
	Tags                Tag        `json:"tags"                       sql:"tags"                       sqlitetype:"text default '{}'"` // A map to store generic info. String and int64 are valid value types of map
	Ignore              int        `json:"-"                          sql:"ignore"                     sqlitetype:"integer"`           // Whether to ignore unit
}

// TableName returns the table which units are stored into.
func (Unit) TableName() string {
	return unitsTableName
}

// TagNames returns a slice of all tag names.
func (u Unit) TagNames(tag string) []string {
	return structset.GetStructFieldTagValues(u, tag)
}

// Usage statistics of each project/tenant/namespace
type Usage struct {
	ID                  int64   `json:"-"                          sql:"id"                         sqlitetype:"integer not null primary key"`
	NumUnits            int64   `json:"num_units"                  sql:"num_units"                  sqlitetype:"integer"` // Number of consumed units
	Project             string  `json:"project"                    sql:"project"                    sqlitetype:"text"`    // Account in batch systems, Tenant in Openstack, Namespace in k8s
	Usr                 string  `json:"usr"                        sql:"usr"                        sqlitetype:"text"`    // Username
	TotalCPUBilling     int64   `json:"total_cpu_billing"          sql:"total_cpu_billing"          sqlitetype:"integer"` // Total CPU billing for project
	TotalGPUBilling     int64   `json:"total_gpu_billing"          sql:"total_gpu_billing"          sqlitetype:"integer"` // Total GPU billing for project
	TotalMiscBilling    int64   `json:"total_misc_billing"         sql:"total_misc_billing"         sqlitetype:"integer"` // Total billing for project that are not in CPU and GPU billing
	AveCPUUsage         float64 `json:"avg_cpu_usage"              sql:"avg_cpu_usage"              sqlitetype:"real"`    // Average CPU usage during lifetime of project
	AveCPUMemUsage      float64 `json:"avg_cpu_mem_usage"          sql:"avg_cpu_mem_usage"          sqlitetype:"real"`    // Average CPU memory during lifetime of project
	TotalCPUEnergyUsage float64 `json:"total_cpu_energy_usage_kwh" sql:"total_cpu_energy_usage_kwh" sqlitetype:"real"`    // Total CPU energy usage in kWh during lifetime of project
	TotalCPUEmissions   float64 `json:"total_cpu_emissions_gms"    sql:"total_cpu_emissions_gms"    sqlitetype:"real"`    // Total CPU emissions in grams during lifetime of project
	AveGPUUsage         float64 `json:"avg_gpu_usage"              sql:"avg_gpu_usage"              sqlitetype:"real"`    // Average GPU usage during lifetime of project
	AveGPUMemUsage      float64 `json:"avg_gpu_mem_usage"          sql:"avg_gpu_mem_usage"          sqlitetype:"real"`    // Average GPU memory during lifetime of project
	TotalGPUEnergyUsage float64 `json:"total_gpu_energy_usage_kwh" sql:"total_gpu_energy_usage_kwh" sqlitetype:"real"`    // Total GPU energy usage in kWh during lifetime of project
	TotalGPUEmissions   float64 `json:"total_gpu_emissions_gms"    sql:"total_gpu_emissions_gms"    sqlitetype:"real"`    // Total GPU emissions in grams during lifetime of project
	TotalIOWriteHot     float64 `json:"total_io_write_hot_gb"      sql:"total_io_write_hot_gb"      sqlitetype:"real"`    // Total IO write on hot storage in GB during lifetime of project
	TotalIOReadHot      float64 `json:"total_io_read_hot_gb"       sql:"total_io_read_hot_gb"       sqlitetype:"real"`    // Total IO read on hot storage in GB during lifetime of project
	TotalIOWriteCold    float64 `json:"total_io_write_cold_gb"     sql:"total_io_write_cold_gb"     sqlitetype:"real"`    // Total IO write on cold storage in GB during lifetime of project
	TotalIOReadCold     float64 `json:"total_io_read_cold_gb"      sql:"total_io_read_cold_gb"      sqlitetype:"real"`    // Total IO read on cold storage in GB during lifetime of project
	TotalIngress        float64 `json:"total_ingress_in_gb"        sql:"total_ingress_in_gb"        sqlitetype:"real"`    // Total ingress traffic in GB of project
	TotalOutgress       float64 `json:"total_outgress_in_gb"       sql:"total_outgress_in_gb"       sqlitetype:"real"`    // Total outgress traffic in GB of project
}

// TableName returns the table which usage stats are stored into.
func (Usage) TableName() string {
	return usageTableName
}

// TagNames returns a slice of all tag names.
func (u Usage) TagNames(tag string) []string {
	return structset.GetStructFieldTagValues(u, tag)
}

// Project struct
type Project struct {
	Name string `json:"name,omitempty" sql:"project" sqlitetype:"text"`
}