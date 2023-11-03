import 'package:flutter/material.dart';
import 'package:intl/intl.dart';
import 'package:swim_results/model/Session.dart';

class SessionView extends StatelessWidget {
  final DateFormat formatter = DateFormat('dd.MM.yyyy');

  final Session s;
  SessionView(this.s, {super.key});

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        Text(formatter.format(s.day)),
      ],

    );
  }
}