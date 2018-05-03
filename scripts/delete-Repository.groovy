def deleteRepository(String name) {

  try {
    repository.getRepositoryManager().delete(name)
  }
  catch (Exception e){
    return e
  }

}

if (args != "") {

  deleteRepository(args)
}

