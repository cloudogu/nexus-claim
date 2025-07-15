def deleteRepository(String name) {

  try {
    repository.@repositoryManager.delete(name)
  }
  catch (Exception e){
    return e
  }
  return null
}

if (args != "") {

  def output = deleteRepository(args)

  return output
}

