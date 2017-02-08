provider "alks" {
    url      = ""
    username = ""
    account  = ""
    role     = ""
    password = ""
}

resource "alks_iamrole" "mah_role" {
    name = "My_Test_Role"
    type = "Amazon EC2"
    include_default_policies = false
}