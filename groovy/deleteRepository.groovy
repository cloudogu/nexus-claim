
def deleteRepository(String name) {
  repository.getRepositoryManager().delete(name)
}

String name

if (args != "") {
  name = args
}

deleteRepository(name)

