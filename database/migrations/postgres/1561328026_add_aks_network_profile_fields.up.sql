ALTER TABLE "azure_aks_clusters" ADD COLUMN "network_plugin" text DEFAULT 'kubenet';
ALTER TABLE "azure_aks_clusters" ADD COLUMN "pod_cidr" varchar(18) DEFAULT '10.244.0.0/16';
ALTER TABLE "azure_aks_clusters" ADD COLUMN "service_cidr" varchar(18) DEFAULT '10.0.0.0/16';
ALTER TABLE "azure_aks_clusters" ADD COLUMN "dns_service_ip" varchar(15) DEFAULT '10.0.0.10';
ALTER TABLE "azure_aks_clusters" ADD COLUMN "docker_bridge_cidr" varchar(18) DEFAULT '172.17.0.1/16';
