ALTER TABLE `azure_aks_clusters` DROP COLUMN `network_plugin`;
ALTER TABLE `azure_aks_clusters` DROP COLUMN `pod_cidr`;
ALTER TABLE `azure_aks_clusters` DROP COLUMN `service_cidr`;
ALTER TABLE `azure_aks_clusters` DROP COLUMN `dns_service_ip`;
ALTER TABLE `azure_aks_clusters` DROP COLUMN `docker_bridge_cidr`;
