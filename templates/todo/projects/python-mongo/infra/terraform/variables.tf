variable "name" {
  type = string
  description = "Name of the the environment which is used to generate a short unique hash used in all resources."
}

variable "location" {
  type = string
  description = "Primary location for all resources"
}

variable "tags" {
  type = "map"
}