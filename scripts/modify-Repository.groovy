import groovy.json.JsonSlurper
import org.sonatype.nexus.repository.Repository

class RepositoryMap {
  Map<String, Map<String, Object>> properties = new HashMap<String, Object>()
}

if (args != "") {
  RepositoryMap repoMap = convertJsonFileToRepoMap(args)
  Repository repo = getRepositoryFromJsonData(repoMap)
  def conf = repo.getConfiguration().copy()
  def attributes = repoMap.properties.get("attributes")
  conf.setAttributes(attributes)

  try {
    repository.getRepositoryManager().update(conf)
  }
  catch (Exception e) {
    return e
  }

  return null
}

Repository getRepositoryFromJsonData(RepositoryMap repoMap){
  return repository.getRepositoryManager().get(repoMap.properties.get("repositoryName"))
}

RepositoryMap convertJsonFileToRepoMap(String jsonData) {
  def inputJson = new JsonSlurper().parseText(jsonData)
  RepositoryMap repo = new RepositoryMap()
  inputJson.each {
    repo.properties.put(it.key, it.value)
  }

  return repo
}
