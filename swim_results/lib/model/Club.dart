class Club {
  
  late int id;
  late String name;
  late String? nationality;

  Club.fromJson(Map json) {
    id = json["id"];
    name = json["name"];
    nationality = json["nationality"];
  }

  Club.fromSwimmerSearchPage(Map json) {
    id = json["clubid"];
    name = json["clubname"];
    nationality = json['nationality'];
  }
}
