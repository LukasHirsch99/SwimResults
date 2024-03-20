import 'package:flutter/material.dart';
import 'package:swim_results/components/api.dart';
import 'package:swim_results/components/drawer.dart';
import 'package:swim_results/components/loadingAnimation.dart';
import 'package:swim_results/event_screen/components/CustomAppBar.dart';
import 'package:swim_results/model/Meet.dart';

import '../components/globals.dart' as globals;
import '../components/meet_item.dart';

class HomePage extends StatefulWidget implements globals.RouteBase {
  static const String routeName = '/home';

  const HomePage({super.key});
  @override
  State<HomePage> createState() => _HomePageState();

  @override
  String getRouteName() {
    return routeName;
  }
}

class _HomePageState extends State<HomePage> with TickerProviderStateMixin {
  int pageIndex = 0;
  late TabController controller;

  @override
  void initState() {
    super.initState();
    controller = TabController(length: 2, vsync: this);
    if (globals.recentMeets.isEmpty || globals.upcomingMeets.isEmpty) {
      update();
    }
  }

  @override
  void dispose() {
    super.dispose();
    controller.dispose();
  }

  @override
  Widget build(BuildContext context) {
    if (globals.upcomingMeets.isEmpty) {
      return const Scaffold(body: LoadingAnimation());
    }

    return Scaffold(
      drawer: const MyDrawer(),
      body: DefaultTabController(
        length: 2,
        child: Builder(
          builder: (context) {
            return NestedScrollView(
              headerSliverBuilder: (context, innerBoxIsScrolled) {
                return [
                  CustomAppBar(
                    controller: controller,
                    tabs: const [
                      Tab(
                        child: Text(
                          "Upcoming",
                          style: TextStyle(fontSize: 18),
                        ),
                      ),
                      Tab(
                        child: Text(
                          "Recent",
                          style: TextStyle(fontSize: 18),
                        ),
                      ),
                    ],
                    title: const Text(
                      "Meets",
                      style: TextStyle(
                        fontSize: 20,
                      ),
                    ),
                  ),
                ];
              },
              body: TabBarView(
                controller: controller,
                children: [
                  eventList(globals.upcomingMeets),
                  eventList(globals.recentMeets)
                ],
              ),
            );
          },
        ),
      ),
    );
  }

  eventList(List<Meet> meets) {
    return CustomScrollView(
      slivers: [
        SliverList(
          delegate: SliverChildBuilderDelegate(
            (context, index) {
              return MeetItem(meets[index]);
            },
            childCount: meets.length,
          ),
        ),
      ],
    );
  }

  Future<void> update() async {
    globals.upcomingMeets = await SwimResultsApi.getUpcomingMeets();
    globals.recentMeets = await SwimResultsApi.getRecentMeets();
    setState(() {});
  }
}
