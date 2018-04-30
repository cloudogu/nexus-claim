def deleteRepository(String name) {
  repository.getRepositoryManager().delete(name)
}

if (args != "") {
  deleteRepository(args)
}



