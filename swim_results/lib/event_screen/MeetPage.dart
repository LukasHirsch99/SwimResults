import 'package:flutter/material.dart';
import 'package:swim_results/components/ColorTheme.dart';
import 'package:swim_results/components/drawer.dart';
import 'package:swim_results/event_screen/components/LiveTimingPage.dart';
import 'package:swim_results/event_screen/components/SwimmerSearchPage.dart';
import 'package:swim_results/event_screen/components/TrainerPage.dart';
import 'package:swim_results/event_screen/components/event_info.dart';
import 'package:swim_results/event_screen/components/SchedulePage.dart';
import 'package:swim_results/model/Meet.dart';
import 'package:swim_results/event_screen/components/SwimmerPage.dart';
import '../components/globals.dart' as globals;

class MeetPage extends StatefulWidget implements globals.RouteBase {
  static const String routeName = '/meetInfo';
  final Meet meet;
  const MeetPage(this.meet, {super.key});

  @override
  State<MeetPage> createState() => _MeetPageState();

  @override
  String getRouteName() {
    return routeName;
  }
}

class _MeetPageState extends State<MeetPage> {
  _MeetPageState();

  bool dataPresent = false;
  int pageIndex = 0;
  PageController pageController = PageController();

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      drawer: const MyDrawer(),
      body: PageView(
        onPageChanged: (page) => setState(() => pageIndex = page),
        controller: pageController,
        allowImplicitScrolling: true,
        children: [
          globals.trainerMode ? TrainerPage(widget.meet) : SwimmerView(globals.myProfile!, widget.meet),
          SwimmerSearchPage(widget.meet),
          SchedulePage(widget.meet),
          MeetOverview(widget.meet),
          LiveTimingPage(widget.meet.id),
        ],
      ),
      bottomNavigationBar: BottomNavigationBar(
        type: BottomNavigationBarType.fixed,
        elevation: 10,
        currentIndex: pageIndex,
        selectedItemColor: colorScheme.primary,
        onTap: (index) {
          setState(() {
            pageController.animateToPage(index,
                duration: const Duration(milliseconds: 250),
                curve: Curves.ease);
            pageIndex = index;
          });
        },
        items: const [
          BottomNavigationBarItem(
            label: 'Me',
            icon: Icon(Icons.person),
          ),
          BottomNavigationBarItem(
            label: 'Search',
            icon: Icon(Icons.search),
          ),
          BottomNavigationBarItem(
            label: "Schedule",
            icon: Icon(Icons.schedule),
          ),
          BottomNavigationBarItem(
            label: "Info",
            icon: Icon(Icons.info),
          ),
          BottomNavigationBarItem(
            icon: Icon(Icons.live_tv_rounded),
            label: "Live",
          )
          // BottomNavigationBarItem(label: 'Favourites', icon: Icon(Icons.group)),
        ],
      ),
    );
  }
}
