import groovy.json.JsonSlurper
import org.sonatype.nexus.blobstore.api.BlobStoreManager
import org.sonatype.nexus.repository.config.Configuration
import org.sonatype.nexus.repository.storage.WritePolicy


class Repository {
    Map<String, Map<String, Object>> properties = new HashMap<String, Object>()
}

def getRepositories(String data) {

    def inputJson = new JsonSlurper().parseText(data)
    Repository repo = new Repository()
    inputJson.each {
        if(it.key.toString().toLowerCase().contains("blobstorename")){
            repo.properties.put(it.key,getBlobStoreName(it.value))
        }
        else if (it.key.toString().toLowerCase().contains("writepolicy")){
            repo.properties.put(it.key,getWritePolicy(it.value))
        }
        else {
            repo.properties.put(it.key, it.value)
        }
    }
 return repo

}

def getBlobStoreName(String value) {
    switch (value.toLowerCase()) {
        case "default" || "" || null: return BlobStoreManager.DEFAULT_BLOBSTORE_NAME
    }
    return value
}

def getWritePolicy(String value){
    switch (value.toLowerCase()){
        case "default" || "allow" || "" || null: return WritePolicy.ALLOW
        case "allow_once" : return WritePolicy.ALLOW_ONCE
        case "deny" : return WritePolicy.DENY
    }
    return value
}


// Work in progress
def createRepository(Repository repo) {

    String name = repo.properties.get("name")
    repo.properties.remove("name")

    String recipeName = repo.properties.get("recipeName")
    repo.properties.remove("recipeName")

    String state = repo.properties.get("_state")
    repo.properties.remove("_state")



    Configuration conf = new Configuration(
            repositoryName: name,
            recipeName: recipeName,
            online: repo.properties.get("online"),
            attributes: [
                    storage: [
                            blobStoreName              : repo.properties.get("blobStoreName"),
                            writePolicy                : writePolicy,
                            strictContentTypeValidation: strictContentTypeValidation
                    ] as Map
            ] as Map
    )

    conf.setAttributes()

    repository.createRepository(conf)

}

if (args != "") {
    def rep = getRepositories(args)
    createRepository(rep)
}


