library msecm.globals;

import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:swim_results/home_screen/home_screen.dart';
import 'package:swim_results/login_screen/login_screen.dart';
import 'package:swim_results/model/Club.dart';
import 'package:swim_results/model/Meet.dart';
import 'package:swim_results/model/Swimmer.dart';
import 'package:swim_results/personal_bests_screen/personal_bests_screen.dart';
import 'package:swim_results/profile_screen/profile_settting_page.dart';

Swimmer? myProfile;
bool trainerMode = false;
Club? myClub;
late SharedPreferences prefs;
List<Meet> upcomingMeets = [];
List<Meet> recentMeets = [];
const primaryColor = Color(0xFF5570FF);
const secondary = Color(0xFF00F2B8);
const defaultBackground = Color(0xFFF8F9FB);

String formatTime(DateTime? t) {
  if (t == null) {
    return "---";
  }
  if (t.minute == 0) {
    return DateFormat("s.S").format(t).cutoffLastDigit();
  } else if (t.hour == 0) {
    return DateFormat("m:ss.S").format(t).cutoffLastDigit();
  }
  return DateFormat("h:mm:ss.S").format(t).cutoffLastDigit();
}

String? getPrefString(String key) => prefs.getString(key);

Swimmer? getProfile() {
  String? swimmerString = prefs.getString("mySwimmer");
  if (swimmerString == null) {
    return null;
  }
  return Swimmer.fromJson(jsonDecode(swimmerString));
}

void setProfile(Swimmer s) {
  myProfile = s;

  prefs.setString(
    "mySwimmer",
    jsonEncode(s.toJson()),
  );
}

Club? getClub() {
  String? clubString = prefs.getString("myClub");
  if (clubString == null) {
    return null;
  }

  return Club.fromJson(jsonDecode(clubString));
}

void setClub(Club c) {
  myClub = c;

  prefs.setString("myClub",
      jsonEncode({"id": c.id, "name": c.name, "nationality": c.nationality}));
}

class Routes {
  static RouteBase profilePage = const ProfileSettingPage();
  static RouteBase homepage = const HomePage();
  static RouteBase recordPage = RecordPage();
  // static RouteBase savedSwimmers = SavedSwimmers();
  // static RouteBase meetInfo = MeetPage();
  static RouteBase loginPage = const LoginPage();
  // static RouteBase startList = StartList();
  // static RouteBase reultsList = ResultList();
  // static RouteBase meetOverview = MeetPage();

  static List<RouteBase> routes = [
    profilePage,
    homepage,
    recordPage,
    loginPage,
  ];

  static Map<String, WidgetBuilder> getRoutes() {
    return {
      for (var item in routes)
        item.getRouteName(): (BuildContext context) => item
    };
  }
}

extension StringExtension on String {
  String capitalize() {
    return "${this[0].toUpperCase()}${substring(1).toLowerCase()}";
  }

  String toTitleCase() {
    return split(' ').map((word) => word.capitalize()).join(' ');
  }

  String toCamelCase() {
    return split(' ').map((word) => word.capitalize()).join(' ');
  }

  String cutoffLastDigit() {
    return substring(0, length - 1);
  }
}

abstract class RouteBase extends Widget {
  const RouteBase({super.key});

  String getRouteName();
}
