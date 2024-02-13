import 'package:flutter/material.dart';
import 'package:flutter_svg/svg.dart';
import 'package:rive/rive.dart' as rive;
import 'package:swim_results/components/ColorTheme.dart';
import 'package:swim_results/components/api.dart';
import 'package:swim_results/components/globals.dart';
import 'package:swim_results/event_screen/components/CustomAppBar.dart';
import 'package:swim_results/event_screen/components/SwimmerPage.dart';
import 'package:swim_results/model/Event.dart';
import 'package:swim_results/model/Heat.dart';
import 'package:swim_results/model/AgeclassResult.dart';
import 'package:swim_results/model/ResultForAgeclass.dart';
import 'package:swim_results/model/Start.dart';

// ignore: must_be_immutable
class EventPage extends StatelessWidget {
  final Event event;
  bool openResultsFirst = false;
  EventPage(
    this.event, {
    super.key,
    this.openResultsFirst = false,
  });

  late List<Heat> heats = [];
  late List<AgeclassResult> ageclassResults = [];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: FutureBuilder(
        future: getData(),
        builder: (context, snapshot) {
          if (!snapshot.hasData || snapshot.data == false) {
            return const Center(
              child: SizedBox(
                width: 200,
                height: 200,
                child: rive.RiveAnimation.asset("assets/loading.riv"),
              ),
            );
          }

          return DefaultTabController(
            initialIndex: openResultsFirst ? 1 : 0,
            length: 2,
            child: NestedScrollView(
              headerSliverBuilder: ((context, innerBoxIsScrolled) {
                return [
                  CustomAppBar(
                    tabs: const [
                      Tab(
                        child: Row(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            Text(
                              "Starts",
                              style: TextStyle(
                                fontSize: 20,
                                fontWeight: FontWeight.w400,
                              ),
                            ),
                            SizedBox(width: 10),
                            Icon(Icons.pool, size: 20),
                          ],
                        ),
                      ),
                      Tab(
                        child: Row(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            Text(
                              "Results",
                              style: TextStyle(
                                fontSize: 20,
                                fontWeight: FontWeight.w400,
                              ),
                            ),
                            SizedBox(width: 10),
                            Icon(Icons.sports_score, size: 20),
                          ],
                        ),
                      ),
                    ],
                    children: [
                      Text(
                        event.name,
                        style: TextStyle(
                          color: colorScheme.onPrimary,
                          fontSize: 20,
                          fontWeight: FontWeight.w300,
                        ),
                      ),
                    ],
                  ),
                ];
              }),
              body: TabBarView(
                // physics: NeverScrollableScrollPhysics(),
                children: [
                  CustomScrollView(
                    slivers: [
                      SliverList(
                        delegate: SliverChildBuilderDelegate(
                          (context, index) {
                            return HeatItem(heats[index]);
                          },
                          childCount: heats.length,
                        ),
                      )
                    ],
                  ),
                  CustomScrollView(
                    slivers: [
                      SliverList(
                        delegate: SliverChildBuilderDelegate(
                          (context, index) {
                            return AgeClassItem(
                              ageclassResults[index],
                            );
                          },
                          childCount: ageclassResults.length,
                        ),
                      ),
                    ],
                  ),
                ],
              ),
            ),
          );
        },
      ),
    );
  }

  Widget AgeClassItem(AgeclassResult r) {
    return Container(
      margin: const EdgeInsets.symmetric(horizontal: 20, vertical: 10),
      decoration: BoxDecoration(
          borderRadius: BorderRadius.circular(10), color: colorScheme.secondary
          // color: const Color(0xFFEBEBEB),
          // boxShadow: const [
          //   BoxShadow(
          //     color: Colors.grey,
          //     blurRadius: 5,
          //     offset: Offset(0, 4),
          //   ),
          // ],
          ),
      child: Column(
        children: [
          Text(
            r.ageClassName,
            style: const TextStyle(
              fontSize: 20,
            ),
          ),
          ListView.builder(
            itemCount: r.results.length,
            shrinkWrap: true,
            padding: const EdgeInsets.all(0),
            physics: const NeverScrollableScrollPhysics(),
            itemBuilder: (context, idx) {
              return Padding(
                padding: const EdgeInsets.symmetric(
                  horizontal: 10,
                  vertical: 2.5,
                ),
                child: ResultItem(r.results[idx], context),
              );
            },
          )
        ],
      ),
    );
  }

  Widget ResultItem(ResultForAgeclass result, BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 15, vertical: 7),
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(10),
        color: colorScheme.primary,
        // boxShadow: const [
        //   BoxShadow(
        //     color: Colors.grey,
        //     blurRadius: 3,
        //     offset: Offset(0, 4),
        //   ),
        // ],
      ),
      child: GestureDetector(
        onTap: () => Navigator.push(
          context,
          MaterialPageRoute(
            builder: (cntx) =>
                SwimmerView(result.result.swimmer!, event.session!.meet!),
          ),
        ),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            Expanded(
              flex: 2,
              child: Row(
                children: [
                  Expanded(
                      child: Icon(
                    Icons.leaderboard_outlined,
                    color: colorScheme.onPrimary,
                  )),
                  Expanded(
                    child: Text(
                      result.ageClass.position != null
                          ? "${result.ageClass.position}."
                          : "---",
                      style: TextStyle(color: colorScheme.onPrimary),
                    ),
                  ),
                ],
              ),
            ),
            Expanded(
              flex: 5,
              child: Text(
                result.result.swimmer!.fullname.toCamelCase(),
                style: TextStyle(color: colorScheme.onPrimary),
              ),
            ),
            Expanded(
              flex: 2,
              child: Text(
                result.result.time != null
                    ? formatTime(result.result.time)
                    : result.result.splits.toString(),
                style: TextStyle(color: colorScheme.onPrimary),
                textAlign: TextAlign.end,
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget HeatItem(Heat heat) {
    return Container(
      margin: const EdgeInsets.symmetric(horizontal: 20, vertical: 10),
      padding: const EdgeInsets.symmetric(vertical: 5),
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(10),
        color: colorScheme.secondary,
        // boxShadow: const [
        //   BoxShadow(
        //     color: Colors.grey,
        //     blurRadius: 5,
        //     offset: Offset(0, 4),
        //   ),
        // ],
      ),
      child: Column(
        children: [
          Text(
            "Heat ${heat.heatNr} / ${heats.length}",
            style: TextStyle(
              color: colorScheme.onSecondary,
              fontSize: 20,
            ),
          ),
          ListView.builder(
            itemCount: heat.starts.length,
            shrinkWrap: true,
            padding: const EdgeInsets.all(0),
            physics: const NeverScrollableScrollPhysics(),
            itemBuilder: (context, idx) {
              return Padding(
                padding: const EdgeInsets.symmetric(
                  horizontal: 10,
                  vertical: 2.5,
                ),
                child: StartItem(heat.starts[idx], context),
              );
            },
          )
        ],
      ),
    );
  }

  Widget StartItem(Start start, BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 5, vertical: 5),
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(10),
        color: colorScheme.primary,
        // boxShadow: const [
        //   BoxShadow(
        //     color: Colors.grey,
        //     blurRadius: 3,
        //     offset: Offset(0, 4),
        //   ),
        // ],
      ),
      child: GestureDetector(
        onTap: () => Navigator.push(
          context,
          MaterialPageRoute(
            builder: (cntx) =>
                SwimmerView(start.swimmer!, event.session!.meet!),
          ),
        ),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            Expanded(
              flex: 2,
              child: Row(
                children: [
                  Expanded(
                    child: SvgPicture.asset(
                      "assets/StartBlock_cleaned.svg",
                      color: colorScheme.onPrimary,
                      height: 22,
                      width: 22,
                    ),
                  ),
                  Expanded(
                    child: Text(
                      "${start.lane}",
                      style: TextStyle(
                        color: colorScheme.onPrimary,
                      ),
                    ),
                  ),
                ],
              ),
            ),
            Expanded(
              flex: 5,
              child: Text(
                start.swimmer!.fullname.toCamelCase(),
                style: TextStyle(
                  color: colorScheme.onPrimary,
                ),
              ),
            ),
            Expanded(
              flex: 2,
              child: Text(
                formatTime(start.time),
                style: TextStyle(
                  color: colorScheme.onPrimary,
                ),
                textAlign: TextAlign.end,
              ),
            ),
          ],
        ),
      ),
    );
  }

  Future<bool> getData() async {
    Future<List<Heat>> heatsFuture = SwimResultsApi.getHeats(event.id);
    Future<List<AgeclassResult>> ageclassResultsFuture =
        SwimResultsApi.getAgeclassResultsForEvent(event.id);

    heats = await heatsFuture;
    ageclassResults = await ageclassResultsFuture;

    return true;
  }
}
