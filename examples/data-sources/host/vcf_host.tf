terraform {
  required_providers {
    vcf = {
      source = "vmware/vcf"
    }
  }
}

provider "vcf" {
  sddc_manager_username = var.sddc_manager_username
  sddc_manager_password = var.sddc_manager_password
  sddc_manager_host     = var.sddc_manager_host
}

data "vcf_host" "example" {
  fqdn = var.host_fqdn
}

output "host_id" {
  value = data.vcf_host.example.id
}

output "host_fqdn" {
  value = data.vcf_host.example.fqdn
}

output "host_hardware" {
  value = [
    for hw in data.vcf_host.example.hardware : {
      hybrid = hw.hybrid
      model  = hw.model
      vendor = hw.vendor
    }
  ]
}

output "host_version" {
  value = data.vcf_host.example.version
}

output "host_status" {
  value = data.vcf_host.example.status
}

output "host_domain" {
  value = [
    for domain in data.vcf_host.example.domain : {
      id = domain.id
    }
  ]
}

output "host_cluster" {
  value = [
    for cluster in data.vcf_host.example.cluster : {
      id = cluster.id
    }
  ]
}

output "host_network_pool" {
  value = [
    for pool in data.vcf_host.example.network_pool : {
      id   = pool.id
      name = pool.name
    }
  ]
}

output "host_cpu" {
  value = [
    for cpu in coalesce(tolist(data.vcf_host.example.cpu), []) : {
      cores             = cpu.cores
      frequency_mhz     = cpu.frequency_mhz
      used_frequency_mhz = cpu.used_frequency_mhz
      cpu_cores         = [
        for core in coalesce(tolist(cpu.cpu_cores), []) : {
          frequency_mhz = core.frequency_mhz
          manufacturer  = core.manufacturer
          model         = core.model
        }
      ]
    }
  ]
}

output "host_memory" {
  value = [
    for mem in data.vcf_host.example.memory : {
      total_capacity_mb = mem.total_capacity_mb
      used_capacity_mb  = mem.used_capacity_mb
    }
  ]
}

output "host_storage" {
  value = [
    for storage in coalesce(tolist(data.vcf_host.example.storage), []) : {
      total_capacity_mb = storage.total_capacity_mb
      used_capacity_mb  = storage.used_capacity_mb
      disks = [
        for disk in coalesce(tolist(storage.disks), []) : {
          capacity_mb   = disk.capacity_mb
          disk_type     = disk.disk_type
          manufacturer  = disk.manufacturer
          model         = disk.model
        }
      ]
    }
  ]
}

output "host_physical_nics" {
  value = [
    for nic in data.vcf_host.example.physical_nics : {
      device_name = nic.device_name
      mac_address = nic.mac_address
      speed       = nic.speed
      unit        = nic.unit
    }
  ]
}

output "host_ip_addresses" {
  value = [
    for ip in data.vcf_host.example.ip_addresses : {
      ip_address = ip.ip_address
      type       = ip.type
    }
  ]
}
