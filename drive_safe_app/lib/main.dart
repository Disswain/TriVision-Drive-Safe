// import 'package:flutter/material.dart';
// import 'screens/home_screen.dart';

// void main() {
//   runApp(const DriveSafeApp());
// }

// class DriveSafeApp extends StatelessWidget {
//   const DriveSafeApp({super.key});

//   @override
//   Widget build(BuildContext context) {
//     return MaterialApp(
//       title: 'Drive Safe',
//       theme: ThemeData(primarySwatch: Colors.blue),
//       home: const HomeScreen(),
//     );
//   }
// }

import 'package:flutter/material.dart';
import 'screens/home_screen.dart';

void main() {
  runApp(const DriveSafeApp());
}

class DriveSafeApp extends StatelessWidget {
  const DriveSafeApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Drive Safe',
      theme: ThemeData(primarySwatch: Colors.blue),
      home: const HomeScreen(),
    );
  }
}
