import 'package:flutter/material.dart';
import 'package:flutter_svg/flutter_svg.dart';
import 'package:swim_results/components/ColorTheme.dart';
import 'package:swim_results/components/globals.dart';
import 'package:swim_results/event_screen/components/EventPage.dart';
import 'package:swim_results/model/Event.dart';

class StartItem extends StatelessWidget {
  final Event event;
  final TextStyle infoStyle = TextStyle(
    fontSize: 20,
    fontWeight: FontWeight.w300,
    color: colorScheme.onPrimary,
  );
  StartItem(this.event, {super.key});

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(10),
        color: colorScheme.primary,
      ),
      child: GestureDetector(
        onTap: () => Navigator.push(
          context,
          MaterialPageRoute(builder: (ctx) => EventPage(event)),
        ),
        child: Column(
          children: [
            Container(
              padding: const EdgeInsets.only(top: 5),
              child: Text(
                event.name,
                style: TextStyle(
                  color: colorScheme.onPrimary,
                  fontSize: 17,
                  fontWeight: FontWeight.w300,
                ),
              ),
            ),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceAround,
              children: [
                Expanded(
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Icon(
                        Icons.format_list_bulleted_rounded,
                        color: colorScheme.onPrimary,
                        size: 25,
                      ),
                      Padding(
                        padding: const EdgeInsets.all(10),
                        child: Text(
                          event.heats[0].heatNr.toString(),
                          style: infoStyle,
                        ),
                      ),
                    ],
                  ),
                ),
                Expanded(
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      SvgPicture.asset(
                        "assets/StartBlock_cleaned.svg",
                        color: colorScheme.onPrimary,
                        height: 22,
                      ),
                      Padding(
                        padding: const EdgeInsets.all(10),
                        child: Text(
                          event.heats[0].starts[0].lane.toString(),
                          style: infoStyle,
                        ),
                      ),
                    ],
                  ),
                ),
                Expanded(
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Icon(
                        Icons.timer_outlined,
                        color: colorScheme.onPrimary,
                        size: 25,
                      ),
                      Padding(
                        padding: const EdgeInsets.all(10),
                        child: Text(
                          formatTime(event.heats[0].starts[0].time),
                          style: infoStyle,
                        ),
                      ),
                    ],
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}
