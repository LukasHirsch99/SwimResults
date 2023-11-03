import 'package:swim_results/model/Session.dart';

class Meet {
  late int id;
  late String name;
  late String? image;
  List<String> invitations = [];
  late DateTime deadline;
  late String address;
  late DateTime startDate;
  late DateTime endDate;
  List<Session> sessions = [];

  Meet.fromJson(Map json) {
    id = json["id"];
    name = json["name"];
    image = json["image"];
    if (json.containsKey("invitations")) {
      for (String invitation in json["invitations"]) {
        invitations.add(invitation);
      }
    }

    deadline = DateTime.parse(json["deadline"]);
    address = json["address"];
    startDate = DateTime.parse(json["startdate"]);
    endDate = DateTime.parse(json["enddate"]);

    if (json.containsKey("session")) {
      for (var session in json["session"]) {
        sessions.add(Session.fromJson(session));
      }
    }
  }
}
