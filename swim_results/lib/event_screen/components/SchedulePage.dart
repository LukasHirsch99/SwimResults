import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import 'package:swim_results/components/ColorTheme.dart';
import 'package:swim_results/components/api.dart';
import 'package:swim_results/components/loadingAnimation.dart';
import 'package:swim_results/event_screen/components/EventPage.dart';
import 'package:swim_results/model/Event.dart';
import 'package:swim_results/model/Meet.dart';
import 'package:swim_results/model/Session.dart';

// ignore: must_be_immutable
class SchedulePage extends StatelessWidget {
  final Meet meet;
  SchedulePage(this.meet, {super.key});

  late List<Session> sessions = [];
  final DateFormat dayFormatter = DateFormat("EEEE dd.MM.yyyy");
  final DateFormat startFormatter = DateFormat("HH:mm");

  @override
  Widget build(BuildContext context) {
    return FutureBuilder(
      future: getData(),
      builder: (context, snapshot) {
        if (!snapshot.hasData) {
          return const LoadingAnimation();
        }

        return ListView.builder(
          itemBuilder: (ctx, sessionIdx) => Session.SessionItemContainer(
              session: sessions[sessionIdx], eventWidget: scheduleEntry),
          itemCount: sessions.length,
        );
      },
    );
  }

  Widget scheduleEntry(BuildContext context, Event event) {
    return GestureDetector(
      onTap: () => {
        if (event.heats.isNotEmpty) {
          Navigator.push(
            context,
            MaterialPageRoute(builder: (cntx) => EventPage(event)),
          ),
        }
      },
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 5, vertical: 7),
        decoration: BoxDecoration(
          borderRadius: BorderRadius.circular(10),
          color: colorScheme.primary,
        ),
        child: Row(
          children: [
            Expanded(
              child: Text(
                event.name,
                style: TextStyle(
                  fontSize: 14,
                  color: colorScheme.onPrimary,
                ),
              ),
            ),
            if (event.heats.isNotEmpty)
              Text(
                "${event.heats.length} Heats",
                style: TextStyle(
                  fontSize: 14,
                  color: colorScheme.onPrimary,
                ),
              ),
          ],
        ),
      ),
    );
  }

  Future<bool> getData() async {
    sessions = await SwimResultsApi.getSessionsForSchedule(meet.id);
    return true;
  }
}
