import 'package:flutter/material.dart';
import 'package:swim_results/components/ColorTheme.dart';

class LiveTimingLane extends StatelessWidget {
  final int? position;
  final int lane;
  final int swimmerId;
  final String swimmerName;
  final String birthyearAndGender;
  final int clubId;
  final String clubName;
  final DateTime? entryTime;
  final DateTime? endTime;
  final String additionalInfo;
  final String splits;

  const LiveTimingLane(
      this.position,
      this.lane,
      this.swimmerId,
      this.swimmerName,
      this.birthyearAndGender,
      this.clubId,
      this.clubName,
      this.entryTime,
      this.endTime,
      this.additionalInfo,
      this.splits,
      {super.key});

  @override
  Widget build(BuildContext context) {
    return Container(
      margin: const EdgeInsets.symmetric(horizontal: 5, vertical: 3),
      padding: const EdgeInsets.symmetric(horizontal: 15, vertical: 7),
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(10),
        color: colorScheme.primary,
        boxShadow: const [
          BoxShadow(
            color: Colors.grey,
            blurRadius: 3,
            offset: Offset(0, 4),
          ),
        ],
      ),
      child: Row(
        children: [
          Text(position.toString()),
          Text(swimmerName.toString()),
        ],
      ),
    );
  }
}

class LiveTiming {
  late String event;
  late int heatNr;
  final List<LiveTimingLane> lanes = [];

  LiveTiming.fromJson(Map<String, dynamic> json) {
    event = json["title1"].split(",")[0];
    heatNr = int.parse(
      RegExp(r"\d+").firstMatch(json["title2"].split(" - ")[1])![0]!,
    );

    for (final lane in json["lanes"]) {
      lanes.add(LiveTimingLane(
        int.parse(RegExp(r"\d+").firstMatch(lane["pl"]!)![0]!),
        int.parse(lane["la"]!),
        int.parse(RegExp(r"\d+").firstMatch(lane["p_ref"]!)![0]!),
        lane["p_name"]!,
        lane["p_detail"]!,
        int.parse(RegExp(r"\d+").firstMatch(lane["c_ref"]!)![0]!),
        lane["c_name"]!,
        DateTime.tryParse("2000-01-01 ${json["etime"]}z"),
        DateTime.tryParse("2000-01-01 ${json["ftime"]}z"),
        lane["add"]!,
        lane["det"]!,
      ));
    }
  }
}
