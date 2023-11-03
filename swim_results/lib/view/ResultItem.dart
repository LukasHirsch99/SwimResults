import 'package:flutter/material.dart';
import 'package:swim_results/components/ColorTheme.dart';
import 'package:swim_results/components/globals.dart';
import 'package:swim_results/event_screen/components/EventPage.dart';
import 'package:swim_results/model/AgeClass.dart';
import 'package:swim_results/model/Result.dart';

class ResultItem extends StatelessWidget {
  final Result result;
  const ResultItem(this.result, {super.key});

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(10),
        color: colorScheme.secondary,
      ),
      child: GestureDetector(
        onTap: () => Navigator.push(
          context,
          MaterialPageRoute(
            builder: (ctx) => EventPage(
              result.event!,
              openResultsFirst: true,
            ),
          ),
        ),
        child: Column(
          children: [
            Padding(
              padding: const EdgeInsets.only(top: 5),
              child: Text(
                result.event!.name,
                style:
                    const TextStyle(fontSize: 17, fontWeight: FontWeight.w300),
              ),
            ),
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                const Icon(
                  Icons.alarm,
                  size: 21,
                ),
                Text(
                  formatTime(result.time),
                  style: const TextStyle(
                    fontSize: 20,
                  ),
                ),
              ],
            ),
            ListView.builder(
              padding: const EdgeInsets.all(0),
              itemCount: result.ageClasses.length,
              shrinkWrap: true,
              physics: const NeverScrollableScrollPhysics(),
              itemBuilder: (ctx, idx) => Padding(
                padding:
                    const EdgeInsets.symmetric(vertical: 2.5, horizontal: 5),
                child: AgeClassItem(result.ageClasses[idx]),
              ),
            ),
            Padding(
              padding: const EdgeInsets.all(5.0),
              child: Text(result.splits!),
            ),
          ],
        ),
      ),
    );
  }

  Widget AgeClassItem(AgeClass ageClass) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 5, vertical: 7),
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
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Padding(
            padding: const EdgeInsets.only(left: 10),
            child: Row(
              children: [
                if (ageClass.position != null)
                  Icon(
                    Icons.leaderboard_outlined,
                    color: colorScheme.onPrimary,
                    size: 21,
                  ),
                const SizedBox(
                  width: 10,
                ),
                Text(
                  ageClass.position != null
                      ? "${ageClass.position}."
                      : ageClass.result!.splits!,
                  style: TextStyle(color: colorScheme.onPrimary),
                ),
              ],
            ),
          ),
          Padding(
            padding: const EdgeInsets.only(right: 10),
            child: Text(
              ageClass.name,
              style: TextStyle(color: colorScheme.onPrimary),
            ),
          ),
        ],
      ),
    );
  }
}
