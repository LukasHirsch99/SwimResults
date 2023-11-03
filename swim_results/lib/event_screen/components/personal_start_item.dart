import 'package:flutter/material.dart';

// import '../../starts_screen/starts_sreen.dart';

class PersonalStartItem extends StatelessWidget {
  final dynamic start;
  final String eventId;

  const PersonalStartItem(this.start, this.eventId, {super.key});

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        const Divider(),
        ListTile(
          onTap: () {
            // Navigator.push(
            //   context,
            //   MaterialPageRoute(
            //     builder: (cntx) =>
            //         StartList(eventId, start['id'], start['name']),
            //   ),
            // );
          },
          title: Text(
            start['name'],
            style: const TextStyle(color: Colors.white, fontSize: 17),
          ),
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
}
