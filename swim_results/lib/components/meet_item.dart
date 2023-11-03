import 'package:flutter/material.dart';
import 'package:swim_results/components/ColorTheme.dart';
import 'package:swim_results/event_screen/MeetPage.dart';
import 'package:intl/intl.dart';
import 'package:swim_results/model/Meet.dart';

class MeetItem extends StatelessWidget {
  final DateFormat startFormat = DateFormat('dd.MM');
  final DateFormat endFormat = DateFormat('dd.MM.yyyy');

  final Meet m;

  MeetItem(this.m, {super.key});

  @override
  Widget build(BuildContext context) {
    String date;
    if (m.startDate == m.endDate) {
      date = startFormat.format(m.startDate);
    } else {
      date =
          "${startFormat.format(m.startDate)}  -  ${startFormat.format(m.endDate)}";
    }

    return Container(
      margin: const EdgeInsets.symmetric(horizontal: 20, vertical: 10),
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(10),
        // color: const Color(0xFFF8F9FB),
        color: colorScheme.primary,
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
          m.name,
          maxLines: 3,
          style: TextStyle(color: colorScheme.onPrimary),
        ),
        subtitle: Text(
          date,
          style: TextStyle(color: colorScheme.onPrimary),
        ),
        leading: thnail(),
        onTap: () {
          Navigator.push(
            context,
            MaterialPageRoute(
              builder: (context) => MeetPage(
                m,
              ),
            ),
          );
        },
      ),
    );
  }

  Widget thnail() {
    if (m.image != null) {
      return Image.network(m.image!, width: 50, height: 50, scale: .4);
    }
    return const Icon(Icons.error_outline, size: 50, color: Colors.grey);
  }
}
