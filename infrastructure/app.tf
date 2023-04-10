resource "fly_app" "sponty-bot" {
  name = "sponty-bot"
  org  = local.fly_io_org
}
