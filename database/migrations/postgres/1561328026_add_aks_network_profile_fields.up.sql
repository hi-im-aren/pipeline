ALTER TABLE "azure_aks_clusters" ADD COLUMN "network_plugin" text;
ALTER TABLE "azure_aks_clusters" ADD COLUMN "pod_cidr" varchar(18);
ALTER TABLE "azure_aks_clusters" ADD COLUMN "service_cidr" varchar(18);
ALTER TABLE "azure_aks_clusters" ADD COLUMN "dns_service_ip" text;
ALTER TABLE "azure_aks_clusters" ADD COLUMN "docker_bridge_cidr" varchar(18);
