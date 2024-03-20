import 'dart:convert';
import 'package:swim_results/components/globals.dart' as globals;
import 'package:http/http.dart';
import 'package:swim_results/model/AgeClass.dart';
import 'package:swim_results/model/Club.dart';
import 'package:swim_results/model/Heat.dart';
import 'package:swim_results/model/LiveTiming.dart';
import 'package:swim_results/model/Meet.dart';
import 'package:swim_results/model/Result.dart';
import 'package:swim_results/model/AgeclassResult.dart';
import 'package:swim_results/model/ResultForAgeclass.dart';
import 'package:swim_results/model/Session.dart' as s;
import 'package:swim_results/model/Swimmer.dart';
import 'package:supabase_flutter/supabase_flutter.dart';

class SwimResultsApi {
  /// the base url for the API
  static const String baseUrl = "pi2b.kq43etgbudgb5lgz.myfritz.net";

  /// Creates an Uri object with given parameters and baseUrl
  static Uri myUrl(String path, [Map<String, String>? queryParameters]) =>
      Uri.http(baseUrl, "v1/$path", queryParameters);

  static Future<List<Meet>> getUpcomingMeets() async {
    final supabase = Supabase.instance.client;
    List data = await supabase
        .from("meet")
        .select('*')
        .gte("enddate", DateTime.now().toIso8601String())
        .order("startdate", ascending: true);

    List<Meet> meets = [];
    for (var meet in data) {
      meets.add(Meet.fromJson(meet));
    }
    return meets;
  }

  static Future<List<Meet>> getRecentMeets() async {
    final supabase = Supabase.instance.client;
    List data = await supabase
        .from("meet")
        .select('*')
        .lt("enddate", DateTime.now().toIso8601String())
        .order("startdate", ascending: false);

    List<Meet> meets = [];
    for (var meet in data) {
      meets.add(Meet.fromJson(meet));
    }
    return meets;
  }

  static Future<Swimmer?> getSwimmerByExactName(String name) async {
    final supabase = Supabase.instance.client;

    List data = await supabase
        .from("swimmer")
        .select('*,club!inner(*)')
        .ilike("firstname", name)
        .ilike("lastname", name);

    if (data.isEmpty) return null;
    Swimmer s = Swimmer.fromJson(data[0]);
    return s;
  }

  /// Get all Results for meet and swimmer
  static Future<List<Result>> getResultsByMeetAndSwimmer(
      int meetId, int swimmerId) async {
    final supabase = Supabase.instance.client;

    List resultData = await supabase
        .from("result")
        .select(
            '*, ageclass!inner(*), swimmer!inner(*, club!inner(*)), event!inner(*, session!inner(*, meet!inner(*)))')
        .eq("swimmerid", swimmerId)
        .eq("event.session.meetid", meetId)
        .order("displaynr", foreignTable: "event", ascending: true);

    if (resultData.isEmpty) return [];

    List<Result> resultList = [];

    for (var result in resultData) {
      resultList.add(Result.fromJson(result));
    }
    return resultList;
  }

  /// Gets all Starts for meet and swimmer
  static Future<List<s.Session>> getStartsByMeetAndSwimmer(
      int meetId, int swimmerId) async {
    final supabase = Supabase.instance.client;

    List startData = await supabase
        .from("session")
        .select(
            '*, meet!inner(*), event!inner(*, heat!inner(*, start!inner(*, swimmer!inner(*, club!inner(*)))))')
        .eq("event.heat.start.swimmerid", swimmerId)
        .eq("meetid", meetId)
        .order("displaynr", foreignTable: "event", ascending: true);

    if (startData.isEmpty) return [];

    List<s.Session> sessionList = [];

    for (var session in startData) {
      sessionList.add(s.Session.fromJson(session));
    }

    return sessionList;
  }

  /// Returns a list of Swimmers filtered by the name
  static Future<List<Swimmer>> getSwimmersByName(String name) async {
    final supabase = Supabase.instance.client;

    List data = await supabase
        .from("swimmer")
        .select("*, club!inner(*)")
        .ilike("lastname", "%$name%");

    if (data.isEmpty) return [];

    List<Swimmer> swimmers = [];

    for (var s in data) {
      swimmers.add(Swimmer.fromJson(s));
    }

    return swimmers;
  }

  /// Returns a list of Clubs filtered by the name
  static Future<List<Club>> getClubsByName(String name) async {
    final supabase = Supabase.instance.client;

    List data =
        await supabase.from("club").select("*").ilike("name", "%$name%");

    if (data.isEmpty) return [];

    List<Club> clubs = [];

    for (var club in data) {
      clubs.add(Club.fromJson(club));
    }

    return clubs;
  }

  static Future<LiveTiming?> getLiveTiming(int meetId) async {
    Uri url = Uri.https(
        "myresults.eu", "/ajax_livetiming.php", {"meet": meetId.toString()});
    Response r = await get(url);
    if (r.statusCode != 200) return null;
    return LiveTiming.fromJson(jsonDecode(r.body));
  }

  static Future<List<Swimmer>> getSwimmersByNameForEvent(
      int meetId, String name) async {
    // Uri url = Uri.https("myresults.eu", "/ajax_searchmeetparticipants.php", {
    //   "language": "en-US",
    //   "meet": meetId.toString(),
    //   "searchstring": name,
    // });
    // Response r = await get(url);
    // List<dynamic> data = jsonDecode(r.body);
    final supabase = Supabase.instance.client;
    // List data = await supabase.from("session").select("meetid, id, event!inner(id, heat!inner(id, start!inner(swimmerid, swimmer!inner(*))))")
    //   .eq("meetid", meetId)
    //   .ilike("event.heat.start.swimmer.name", "%$name%");
    List data = await supabase.rpc("getswimmersbynameformeet",
        params: {"meetingid": meetId, "swimmername": "$name%"});

    return List.generate(
        data.length, (i) => Swimmer.fromSwimmerSearchPage(data[i]));
  }

  /// Updates the database if new events are available
  static Future<bool> updateDatabase() async {
    await get(myUrl("update"));
    return true;
  }

  /// Gets the personal best times of the user
  static Future<Map> getRecords() async {
    Uri url = SwimResultsApi.myUrl("records", {
      "lastname": globals.myProfile!.fullname,
    });
    Response r = await get(url);
    if (r.body == 'ERROR') {
      print('Error fetching records');
      return {};
    }
    return jsonDecode(r.body);
  }

  /// Gets the starts for a event (eventId) and discipline (startId)
  static Future<List<Heat>> getHeats(int eventId) async {
    final supabase = Supabase.instance.client;
    List data = await supabase
        .from("heat")
        .select("*, start!inner(*, swimmer!inner(*, club!inner(*)))")
        .eq("eventid", eventId);

    List<Heat> heats = [];

    for (var heat in data) {
      heats.add(Heat.fromJson(heat));
    }
    return heats;
  }

  static Future<List<AgeclassResult>> getAgeclassResultsForEvent(
      int eventId) async {
    final supabase = Supabase.instance.client;
    List data = await supabase
        .from("result")
        .select(
            "*, ageclass!inner(*), swimmer!inner(*, club!inner(*)), event!inner(*, session!inner(*, meet!inner(*)))")
        .eq("eventid", eventId);

    List<AgeclassResult> resultsForEvent = [];
    Map<String, List<ResultForAgeclass>> m = {};

    for (Map result in data) {
      for (Map ac in result["ageclass"]) {
        m.putIfAbsent(ac["name"], () => []).add(
            ResultForAgeclass(Result.fromJson(result), AgeClass.fromJson(ac)));
      }
    }

    m.forEach((key, value) {
      resultsForEvent.add(AgeclassResult(key, value));
    });

    return resultsForEvent;
  }

  static Future<List<s.Session>> getSessionsForSchedule(int meetId) async {
    final supabase = Supabase.instance.client;

    List sessions = await supabase
        .from("session")
        .select('*, event!inner(*, heat(*)), meet!inner(*)')
        .eq("meetid", meetId)
        .order("displaynr", foreignTable: "event", ascending: true);

    if (sessions.isEmpty) return [];

    List<s.Session> sessionList = [];

    for (var session in sessions) {
      sessionList.add(s.Session.fromJson(session));
    }

    return sessionList;
  }

  static Future<List<dynamic>> getScheduleOld(dynamic eventId) async {
    Uri url = myUrl("schedule", {"eventId": eventId.toString()});
    Response r = await get(url);
    return jsonDecode(r.body).toList();
  }
}
