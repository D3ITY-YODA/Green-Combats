package com.greencompass.feature.main

import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.PaddingValues
import androidx.compose.foundation.layout.padding
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Explore
import androidx.compose.material.icons.filled.Home
import androidx.compose.material.icons.filled.Notifications
import androidx.compose.material.icons.filled.Person
import androidx.compose.material.icons.filled.Edit
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import com.greencompass.navigation.AppRoute

@Composable
fun MainAppScaffold(
    currentRoute: String?,
    onNavigate: (AppRoute) -> Unit,
    content: @Composable (PaddingValues) -> Unit
) {
    Scaffold(
        bottomBar = {
            NavigationBar(containerColor = Color.White) {
                NavigationBarItem(selected = currentRoute == "today", onClick = { onNavigate(AppRoute.Today) }, icon = { Icon(Icons.Default.Home, contentDescription = "Today") }, label = { Text("Today") })
                NavigationBarItem(selected = currentRoute == "explore", onClick = { onNavigate(AppRoute.Explore) }, icon = { Icon(Icons.Default.Explore, contentDescription = "Explore") }, label = { Text("Explore") })
                NavigationBarItem(selected = currentRoute == "updates", onClick = { onNavigate(AppRoute.Updates) }, icon = { Icon(Icons.Default.Notifications, contentDescription = "Updates") }, label = { Text("Updates") })
                NavigationBarItem(selected = currentRoute == "report", onClick = { onNavigate(AppRoute.Report) }, icon = { Icon(Icons.Default.Edit, contentDescription = "Report") }, label = { Text("Report") })
                NavigationBarItem(selected = currentRoute == "profile", onClick = { onNavigate(AppRoute.Profile) }, icon = { Icon(Icons.Default.Person, contentDescription = "Profile") }, label = { Text("Profile") })
            }
        }
    ) { innerPadding ->
        Box(modifier = Modifier.padding(innerPadding)) { content(innerPadding) }
    }
}
