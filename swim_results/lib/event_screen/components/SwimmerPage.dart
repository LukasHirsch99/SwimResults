import 'package:flutter/material.dart';
import 'package:intl/date_symbol_data_local.dart';
import 'package:intl/intl.dart';
import 'package:swim_results/components/ColorTheme.dart';
import 'package:swim_results/components/api.dart';
import 'package:swim_results/components/drawer.dart';
import 'package:swim_results/components/loadingAnimation.dart';
import 'package:swim_results/event_screen/components/CustomAppBar.dart';
import 'package:swim_results/model/Event.dart';
import 'package:swim_results/model/Result.dart';
import 'package:swim_results/view/ResultItem.dart';
import 'package:swim_results/components/globals.dart';
import 'package:swim_results/model/Meet.dart';
import 'package:swim_results/model/Session.dart';
import 'package:swim_results/model/Swimmer.dart';

class SwimmerView extends StatefulWidget {
  final Swimmer swimmer;
  final Meet m;

  const SwimmerView(this.swimmer, this.m, {super.key});

  @override
  State<SwimmerView> createState() => _SwimmerViewState();
}

class _SwimmerViewState extends State<SwimmerView>
    with TickerProviderStateMixin {
  late List<Session> sessions;
  late List<Result> results;
  bool dataPresent = false;

  final DateFormat dayFormatter = DateFormat("EEEE dd.MM.yyyy");
  final DateFormat startFormatter = DateFormat("HH:mm");
  late TabController tabController;

  Future<bool> retTrue() async => true;

  @override
  void initState() {
    super.initState();
    Intl.defaultLocale = "de_DE";
    initializeDateFormatting("de_DE");
    tabController = TabController(length: 2, vsync: this);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      drawer: const MyDrawer(),
      body: FutureBuilder(
        future: dataPresent ? retTrue() : getData(),
        builder: (context, snapshot) {
          if (snapshot.hasData && snapshot.data == true) {
            return DefaultTabController(
              length: sessions.isNotEmpty ? 2 : 0,
              child: Builder(
                builder: (context) {
                  return NestedScrollView(
                      headerSliverBuilder: ((context, innerBoxIsScrolled) {
                        return [
                          CustomAppBar(
                            tabs: sessions.isNotEmpty
                                ? const [
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
                                  ]
                                : [],
                            children: [
                              Text(
                                widget.swimmer.fullname,
                                style: TextStyle(
                                  color: colorScheme.onPrimary,
                                  fontSize: 20,
                                  fontWeight: FontWeight.w300,
                                ),
                              ),
                              Text(
                                widget.swimmer.club!.name,
                                style: TextStyle(
                                  color: colorScheme.onPrimary,
                                  fontSize: 18,
                                  fontWeight: FontWeight.w300,
                                ),
                              ),
                              Text(
                                "${widget.swimmer.birthyear} - ${widget.swimmer.gender.capitalize()}",
                                style: TextStyle(
                                  color: colorScheme.onPrimary,
                                  fontSize: 18,
                                  fontWeight: FontWeight.w300,
                                ),
                              ),
                            ],
                          ),
                        ];
                      }),
                      body: sessions.isNotEmpty
                          ? TabBarView(
                              children: [
                                CustomScrollView(
                                  shrinkWrap: true,
                                  slivers: [
                                    SliverList(
                                      delegate: SliverChildBuilderDelegate(
                                        (context, sessionIdx) {
                                          return Session.SessionItemContainer(
                                            session: sessions[sessionIdx],
                                            eventWidget:
                                                Event.StartItemForSwimmerPage,
                                          );
                                        },
                                        childCount: sessions.length,
                                      ),
                                    )
                                  ],
                                ),
                                CustomScrollView(
                                  shrinkWrap: true,
                                  slivers: [
                                    SliverList(
                                      delegate: SliverChildBuilderDelegate(
                                        (context, index) {
                                          return Padding(
                                            padding: const EdgeInsets.only(
                                              top: 10,
                                              left: 10,
                                              right: 10,
                                            ),
                                            child: ResultItem(
                                              results[index],
                                            ),
                                          );
                                        },
                                        childCount: results.length,
                                      ),
                                    ),
                                  ],
                                ),
                              ],
                            )
                          : const Center(
                              child: SizedBox(
                                width: 200,
                                height: 200,
                                child: Text(
                                  "Looks pretty empty in here",
                                  textAlign: TextAlign.center,
                                  style: TextStyle(
                                    fontSize: 20,
                                  ),
                                ),
                              ),
                            ));
                },
              ),
            );
          }
          return const LoadingAnimation();
        },
      ),
    );
  }

  Future<bool> getData() async {
    Future<List<Session>> sessionsFuture =
        SwimResultsApi.getStartsByMeetAndSwimmer(
            widget.m.id, widget.swimmer.id);

    Future<List<Result>> resultsFuture =
        SwimResultsApi.getResultsByMeetAndSwimmer(
            widget.m.id, widget.swimmer.id);

    sessions = await sessionsFuture;
    results = await resultsFuture;

    setState(() => dataPresent = true);
    return true;
  }
}
