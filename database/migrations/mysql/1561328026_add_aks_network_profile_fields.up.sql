ALTER TABLE `azure_aks_clusters` ADD COLUMN `network_plugin`;
ALTER TABLE `azure_aks_clusters` ADD COLUMN `pod_cidr`           varchar(18) DEFAULT NULL;
ALTER TABLE `azure_aks_clusters` ADD COLUMN `service_cidr`       varchar(18) DEFAULT NULL;
ALTER TABLE `azure_aks_clusters` ADD COLUMN `dns_service_ip`;
ALTER TABLE `azure_aks_clusters` ADD COLUMN `docker_bridge_cidr` varchar(18) DEFAULT NULL;
