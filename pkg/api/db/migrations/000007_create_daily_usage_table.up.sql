CREATE TABLE IF NOT EXISTS daily_usage (
 "id" integer not null primary key,
 "resource_manager" text default "",
 "cluster_id" text, 
 "num_units" integer,
 "project" text,
 "groupname" text,
 "username" text,
 "total_time_seconds" text default '{}', 
 "avg_cpu_usage" text default '{}', 
 "avg_cpu_mem_usage" text default '{}',
 "total_cpu_energy_usage_kwh" text default '{}', 
 "total_cpu_emissions_gms" text default '{}',
 "avg_gpu_usage" text default '{}', 
 "avg_gpu_mem_usage" text default '{}',
 "total_gpu_energy_usage_kwh" text default '{}', 
 "total_gpu_emissions_gms" text default '{}',
 "total_io_write_stats" text default '{}', 
 "total_io_read_stats" text default '{}',
 "total_ingress_stats" text default '{}', 
 "total_outgress_stats" text default '{}',
 "num_updates" integer default 0,  
 "last_updated_at" text
);
CREATE UNIQUE INDEX uq_cluster_id_project_usr_lastupdated ON daily_usage (cluster_id,username,project,last_updated_at);
