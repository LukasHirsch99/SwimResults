import 'package:swim_results/model/ResultForAgeclass.dart';

class AgeclassResult {
  final String ageClassName;
  final List<ResultForAgeclass> results;

  AgeclassResult(this.ageClassName, this.results) {
    results.sort((a, b) {
      if (a.ageClass.position == null && b.ageClass.position == null) {
        return 0;
      }

      if (b.ageClass.position == null) {
        return -1;
      }

      if (a.ageClass.position == null) {
        return 1;
      }

      if (a.ageClass.position! < b.ageClass.position!) {
        return -1;
      }

      return 1;
    });
  }
}
