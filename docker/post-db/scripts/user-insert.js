function get_results(result) {
    print(tojson(result));
}

function insert_post(object) {
    print(db.posts.insert(object));
}

insert_post({
    "_id": ObjectId("57a98d98e4b00679b4a830af"),
    "username": "yagiyu",
    "password": "be9cc41f942ff8ebc3a24f63240a9c215e6bcf5a",
    "email": "yagiyu@gmail.com",
    "salt": "c748112bc027878aa62812ba1ae00e40ad46d497"
});

print("_______POST DATA_______");
db.post.find().forEach(get_results);
print("______END POST DATA_____");
