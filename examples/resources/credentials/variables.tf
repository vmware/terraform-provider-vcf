variable "esxi_host_1" {
  default     = "10.10.01.01"
  description = "Name of the first host in the setup"
}

variable "esxi_host_1_pass" {
  default     = "s0m3_p@$w0rd"
  description = "The new password for the host"
  sensitive   = true
}