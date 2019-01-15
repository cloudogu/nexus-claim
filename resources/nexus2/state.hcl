repository "present-repo" {
  Name = "Repository which is present"
  _state = "present"
}

repository "absent-repo" {
  Name = "Repository which is absent"
  _state = "absent"
}
