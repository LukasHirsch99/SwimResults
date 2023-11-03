import 'package:flutter/material.dart';
import 'package:swim_results/components/ColorTheme.dart';
import 'package:swim_results/components/api.dart';
import 'package:swim_results/components/drawer.dart';
import 'package:swim_results/components/loadingAnimation.dart';
import 'package:swim_results/event_screen/components/CustomAppBar.dart';
import 'package:swim_results/model/Meet.dart';
import 'package:rive/rive.dart' as rive;
import '../components/meet_item.dart';
import '../components/globals.dart' as globals;

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
  PageController pageController = PageController();

  @override
  void initState() {
    super.initState();
    if (globals.recentMeets.isEmpty || globals.upcomingMeets.isEmpty) {
      update();
    }
  }

  @override
  Widget build(BuildContext context) {
    if (globals.upcomingMeets.isEmpty) {
      return const Scaffold(
        body: LoadingAnimation()
      );
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
                    tabs: const [
                      Tab(
                        child: Row(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            Text(
                              "Upcoming",
                              style: TextStyle(fontSize: 20),
                            ),
                            SizedBox(width: 10),
                            Icon(Icons.event_note, size: 20),
                          ],
                        ),
                      ),
                      Tab(
                        child: Row(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            Text(
                              "Recent",
                              style: TextStyle(fontSize: 20),
                            ),
                            SizedBox(width: 10),
                            Icon(Icons.event_repeat_outlined, size: 20),
                          ],
                        ),
                      ),
                    ],
                    children: [
                      Text(
                        "Meets",
                        style: TextStyle(
                          color: colorScheme.onBackground,
                          fontSize: 20,
                        ),
                      ),
                    ],
                  ),
                ];
              },
              body: TabBarView(
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

  eventList(List<Meet> e) {
    return CustomScrollView(
      slivers: [
        SliverList(
          delegate: SliverChildBuilderDelegate(
            (context, index) {
              return MeetItem(e[index]);
            },
            childCount: e.length,
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
