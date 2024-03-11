import 'package:swim_results/components/globals.dart';
import 'package:swim_results/model/Club.dart';

class Swimmer {
  late int id;
  late String fullname;
  late String firstname;
  late String lastname;
  int? birthyear;
  late int clubId;
  Club? club;
  late String gender;
  bool? isrelay;

  Swimmer.fromJson(Map json) {
    id = json["id"];
    firstname = json["firstname"];
    lastname = json["lastname"];
    fullname = "$lastname $firstname";
    birthyear = json["birthyear"];
    clubId = json["clubid"];
    if (json.containsKey("club")) club = Club.fromJson(json["club"]);
    gender = json["gender"];
    isrelay = json["isrelay"];
  }

  Swimmer.fromSwimmerSearchPage(Map json) {
    id = json["id"];
    firstname = json["firstname"];
    lastname = json["lastname"];
    fullname = "$lastname $firstname";
    birthyear = json["birthyear"];
    clubId = json["clubid"];
    gender = json["gender"];
    club = Club.fromSwimmerSearchPage(json);
    isrelay = json["isrelay"];
  }

  Map<String, dynamic> toJson() => {
      "id": id,
      "name": fullname,
      "birthyear": birthyear,
      "clubid": clubId,
      "gender": gender,
      "firstname": firstname,
      "lastname": lastname,
      "club": {
        "id": club?.id,
        "name": club?.name,
        "nationality": club?.nationality
      },
  };
}
