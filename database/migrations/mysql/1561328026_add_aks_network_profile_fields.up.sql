ALTER TABLE `azure_aks_clusters` ADD COLUMN `network_plugin`     varchar(255) COLLATE 'utf8mb4_unicode_ci' NULL DEFAULT 'kubenet';
ALTER TABLE `azure_aks_clusters` ADD COLUMN `pod_cidr`           varchar(18) COLLATE 'utf8mb4_unicode_ci' NULL DEFAULT '10.244.0.0/16';
ALTER TABLE `azure_aks_clusters` ADD COLUMN `service_cidr`       varchar(18) COLLATE 'utf8mb4_unicode_ci' NULL DEFAULT '10.0.0.0/16';
ALTER TABLE `azure_aks_clusters` ADD COLUMN `dns_service_ip`     varchar(15) COLLATE 'utf8mb4_unicode_ci' NULL DEFAULT '10.0.0.10';
ALTER TABLE `azure_aks_clusters` ADD COLUMN `docker_bridge_cidr` varchar(18) COLLATE 'utf8mb4_unicode_ci' NULL DEFAULT '172.17.0.1/16';

UPDATE `azure_aks_clusters` SET network_plugin     = 'kubenet';
UPDATE `azure_aks_clusters` SET pod_cidr           = '10.244.0.0/16';
UPDATE `azure_aks_clusters` SET service_cidr       = '10.0.0.0/16';
UPDATE `azure_aks_clusters` SET dns_service_ip     = '10.0.0.10';
UPDATE `azure_aks_clusters` SET docker_bridge_cidr = '172.17.0.1/16';
