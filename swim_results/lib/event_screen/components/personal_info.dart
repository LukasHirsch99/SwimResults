import 'package:flutter/material.dart';
import 'package:swim_results/event_screen/components/personal_result_item.dart';
import 'package:swim_results/model/Meet.dart';
import 'package:swim_results/model/Session.dart';
import 'package:swim_results/view/SessionView.dart';

// import '../../result_screen/results_screen.dart';
// import '../../starts_screen/starts_sreen.dart';

class PersonalInfo extends StatelessWidget {
  final Future<void> Function() onRefresh;
  final Meet m;
  final String name;
  final List<Session> personalStarts, personalResults;

  const PersonalInfo(this.onRefresh, this.personalStarts, this.personalResults,
      this.m, this.name, {super.key});

  @override
  Widget build(BuildContext context) {
    return RefreshIndicator(
      onRefresh: onRefresh,
      child: SingleChildScrollView(
        child: Column(
          children: [
            const SizedBox(
              height: 10,
            ),
            Container(
              margin: const EdgeInsets.symmetric(horizontal: 5),
              padding: const EdgeInsets.symmetric(vertical: 5),
              alignment: Alignment.center,
              decoration: BoxDecoration(
                  borderRadius: BorderRadius.circular(100),
                  color: Colors.tealAccent[400]),
              child: const Text('Starts',
                  style: TextStyle(fontWeight: FontWeight.bold, fontSize: 20)),
            ),
            for (var s in personalStarts)
              // PersonalStartItem(start: s, eventId: m.id.toString()),
              SessionView(s),
            Container(
              padding: const EdgeInsets.symmetric(vertical: 5),
              margin: const EdgeInsets.symmetric(horizontal: 5),
              alignment: Alignment.center,
              decoration: BoxDecoration(
                  borderRadius: BorderRadius.circular(100),
                  color: Colors.tealAccent[400]),
              child: const Text('Results',
                  style: TextStyle(fontWeight: FontWeight.bold, fontSize: 20)),
            ),
            for (var r in personalResults)
              PersonalResultItem(r, m.id.toString())
          ],
        ),
      ),
    );
  }

  start(context, start) {
    return Column(
      children: [
        const Divider(),
        ListTile(
          onTap: () {
            // Navigator.push(
            //   context,
            //   MaterialPageRoute(
            //     builder: (cntx) => StartList(m.id, start['id'], start['name']),
            //   ),
            // );
          },
          title: Text(start['name'],
              style: const TextStyle(color: Colors.white, fontSize: 17)),
          subtitle: Text(
            start['heat'],
            style: const TextStyle(color: Colors.white),
          ),
          leading: Padding(
            padding: const EdgeInsets.symmetric(vertical: 16),
            child: Text(start['time'],
                style: const TextStyle(fontWeight: FontWeight.bold)),
          ),
        ),
      ],
    );
  }

  result(context, result) {
    TextStyle style;
    if (result['place'] == '3.') {
      style = TextStyle(
          color: Colors.deepOrange[300],
          fontWeight: FontWeight.bold,
          fontSize: 25);
    } else if (result['place'] == '2.')
      style = const TextStyle(
          color: Colors.blueGrey, fontWeight: FontWeight.bold, fontSize: 25);
    else if (result['place'] == '1.')
      style = const TextStyle(
          color: Colors.yellow, fontWeight: FontWeight.bold, fontSize: 25);
    else
      style = const TextStyle(fontWeight: FontWeight.bold, fontSize: 25);

    return Column(
      children: [
        const Divider(),
        ListTile(
          onTap: () {
            // Navigator.push(
            //   context,
            //   MaterialPageRoute(
            //     builder: (cntx) =>
            //         ResultList(m.id, result['id'], result['name']),
            //   ),
            // );
          },
          title: Text(result['name']),
          trailing: result['place'] != ''
              ? Text(result['place'], style: style)
              : const Text('DQ', style: TextStyle(fontWeight: FontWeight.bold)),
          leading: Padding(
            padding: const EdgeInsets.symmetric(vertical: 16),
            child: Text(
              result['time'],
              style: const TextStyle(fontWeight: FontWeight.bold),
            ),
          ),
          subtitle: Text('${result['timeInfo']}\n${result['class']}',
              style: const TextStyle(color: Colors.white)),
        ),
      ],
    );
  }
}
