import 'package:flutter/material.dart';
// import 'package:swim_results/result_screen/results_screen.dart';

// ignore: must_be_immutable
class PersonalResultItem extends StatelessWidget {
  final dynamic result;
  final String eventId;
  late TextStyle style;

  PersonalResultItem(this.result, this.eventId, {super.key}) {
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
  }

  @override
  Widget build(BuildContext context) {
    return ListTile(
      onTap: () {
        // Navigator.push(
        //   context,
        //   MaterialPageRoute(
        //     builder: (cntx) =>
        //         ResultList(eventId, result['id'], result['name']),
        //   ),
        // );
      },
      title: Text(result['name']),
      leading: result['place'] != ''
          ? Text(result['place'], style: style)
          : const Text('DQ', style: TextStyle(fontWeight: FontWeight.bold)),
      trailing: Padding(
        padding: const EdgeInsets.symmetric(vertical: 16),
        child: Text(
          result['time'],
          style: const TextStyle(fontWeight: FontWeight.bold),
        ),
      ),
      subtitle: Text('${result['timeInfo']}\n${result['class']}',
          style: const TextStyle(color: Colors.white)),
    );
  }
}
