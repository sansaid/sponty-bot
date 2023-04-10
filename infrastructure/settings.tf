terraform {
  cloud {
    organization = "sanyia"

    workspaces {
      name = "sponty-bot"
    }
  }

  required_providers {
    fly = {
      source  = "fly-apps/fly"
      version = "~> 0.0.20"
    }
  }
}
