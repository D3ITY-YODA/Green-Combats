package com.greencompass.feature.main

import androidx.compose.foundation.layout.*
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Edit
import androidx.compose.material.icons.filled.Home
import androidx.compose.material.icons.filled.Notifications
import androidx.compose.material.icons.filled.Person
import androidx.compose.material.icons.filled.Search
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.unit.dp
import androidx.navigation.NavController
import androidx.navigation.NavGraph.Companion.findStartDestination
import androidx.navigation.compose.currentBackStackEntryAsState
import com.greencompass.core.ui.GreenCompassColors
import com.greencompass.core.ui.GreenCompassTypography
import com.greencompass.navigation.AppRoute

data class BottomNavItem(val route: AppRoute, val label: String, val icon: ImageVector)

val bottomNavItems = listOf(
    BottomNavItem(AppRoute.Today, "Today", Icons.Filled.Home),
    BottomNavItem(AppRoute.Explore, "Explore", Icons.Filled.Search),
    BottomNavItem(AppRoute.Updates, "Updates", Icons.Filled.Notifications),
    BottomNavItem(AppRoute.Report, "Report", Icons.Filled.Edit),
    BottomNavItem(AppRoute.Profile, "Profile", Icons.Filled.Person)
)

@Composable
fun MainAppScaffold(
    navController: NavController,
    content: @Composable () -> Unit
) {
    val navBackStackEntry by navController.currentBackStackEntryAsState()
    // Get the fully qualified route name (e.g., "com.greencompass.navigation.AppRoute.Today")
    val currentRoute = navBackStackEntry?.destination?.route ?: ""
    
    // Check if the current route matches the qualified name of any bottom nav item
    val showBottomNav = bottomNavItems.any { item -> 
        currentRoute == item.route::class.qualifiedName 
    }

    Scaffold(
        bottomBar = {
            if (showBottomNav) {
                NavigationBar(
                    containerColor = Color.White,
                    contentColor = GreenCompassColors.Charcoal,
                    tonalElevation = 0.dp
                ) {
                    bottomNavItems.forEach { item ->
                        // Check if this specific item is the active one
                        val selected = currentRoute == item.route::class.qualifiedName
                        
                        NavigationBarItem(
                            selected = selected,
                            onClick = {
                                navController.navigate(item.route) {
                                    popUpTo(navController.graph.findStartDestination().id) {
                                        saveState = true
                                    }
                                    launchSingleTop = true
                                    restoreState = true
                                }
                            },
                            icon = {
                                Icon(
                                    imageVector = item.icon,
                                    contentDescription = item.label,
                                    tint = if (selected) GreenCompassColors.ForestGreen else GreenCompassColors.MutedText
                                )
                            },
                            label = {
                                Text(
                                    text = item.label,
                                    style = GreenCompassTypography.labelMedium,
                                    color = if (selected) GreenCompassColors.ForestGreen else GreenCompassColors.MutedText
                                )
                            },
                            colors = NavigationBarItemDefaults.colors(
                                selectedIconColor = GreenCompassColors.ForestGreen,
                                unselectedIconColor = GreenCompassColors.MutedText,
                                selectedTextColor = GreenCompassColors.ForestGreen,
                                unselectedTextColor = GreenCompassColors.MutedText,
                                indicatorColor = GreenCompassColors.SoftSage
                            )
                        )
                    }
                }
            }
        },
        containerColor = Color.White,
        contentColor = GreenCompassColors.Charcoal
    ) { paddingValues ->
        Box(modifier = Modifier.padding(paddingValues)) {
            content()
        }
    }
}
