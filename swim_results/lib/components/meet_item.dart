import 'package:flutter/material.dart';
import 'package:swim_results/components/ColorTheme.dart';
import 'package:swim_results/event_screen/MeetPage.dart';
import 'package:intl/intl.dart';
import 'package:swim_results/model/Meet.dart';

class MeetItem extends StatelessWidget {
  final DateFormat startFormat = DateFormat('dd.MM');
  final DateFormat endFormat = DateFormat('dd.MM.yyyy');

  final Meet meet;

  MeetItem(this.meet, {super.key});

  @override
  Widget build(BuildContext context) {
    String date;
    if (meet.startDate == meet.endDate) {
      date = startFormat.format(meet.startDate);
    } else {
      date =
          "${startFormat.format(meet.startDate)}  -  ${startFormat.format(meet.endDate)}";
    }

    return Container(
      margin: const EdgeInsets.only(left: 10, right: 10, top: 20),
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(10),
        // color: const Color(0xFFF8F9FB),
        color: colorScheme.surfaceTint,
        // boxShadow: [
        //   BoxShadow(
        //     color: colorScheme.onPrimary,
        //     blurRadius: 10,
        //     offset: Offset(3, 3),
        //   ),
        // ],
      ),
      child: ListTile(
        title: Text(
          meet.name,
          style: const TextStyle(fontSize: 15),
          maxLines: 3,
        ),
        subtitle: Column(
          children: [
            const SizedBox(height: 5),
            Row(
              children: [
                Icon(
                  Icons.event_rounded,
                  size: 18,
                  color: colorScheme.primary,
                ),
                const SizedBox(width: 10),
                Text(
                  date,
                ),
              ],
            ),
            const SizedBox(height: 5),
            Row(
              children: [
                Icon(
                  Icons.location_on_rounded,
                  size: 18,
                  color: colorScheme.primary,
                ),
                const SizedBox(width: 10),
                Expanded(
                  child: Text(
                    meet.address,
                    overflow: TextOverflow.ellipsis,
                  ),
                ),
              ],
            ),
          ],
        ),
        // leading: thnail(),
        onTap: () {
          Navigator.push(
            context,
            MaterialPageRoute(
              builder: (context) => MeetPage(
                meet,
              ),
            ),
          );
        },
      ),
    );
  }

  Widget thnail() {
    if (meet.image != null) {
      return Image.network(meet.image!, width: 50, height: 50, scale: .4);
    }
    return const Icon(Icons.error_outline, size: 50, color: Colors.grey);
  }
}
