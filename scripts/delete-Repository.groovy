def deleteRepository(String name) {

  try {
    repository.getRepositoryManager().delete(name)
  }
  catch (Exception e){
    return e
  }

  return "successfully deleted " + name

}

if (args != "") {

  def output = deleteRepository(args)

  return output
}

