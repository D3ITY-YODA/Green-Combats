package com.greencompass.navigation

import androidx.compose.runtime.Composable
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.rememberNavController
import androidx.navigation.toRoute
import com.greencompass.feature.explore.ExploreScreen
import com.greencompass.feature.main.MainAppScaffold
import com.greencompass.feature.onboarding.*
import com.greencompass.feature.places.SavedPlacesScreen
import com.greencompass.feature.profile.ProfileScreen
import com.greencompass.feature.reports.ReportFormScreen
import com.greencompass.feature.reports.ReportScreen
import com.greencompass.feature.reports.ReportSuccessScreen
import com.greencompass.feature.today.TodayRoute
import com.greencompass.feature.updates.UpdateDetailScreen
import com.greencompass.feature.updates.UpdatesScreen

@Composable
fun AppNavHost() {
    val navController = rememberNavController()

    NavHost(navController = navController, startDestination = AppRoute.Welcome) {
        // Onboarding
        composable<AppRoute.Welcome> { WelcomeScreen(onGetStarted = { navController.navigate(AppRoute.AccountChoice) }, onChooseLanguage = { navController.navigate(AppRoute.LanguageSelection) }) }
        composable<AppRoute.LanguageSelection> { LanguageSelectionScreen(onContinue = { navController.navigate(AppRoute.AccountChoice) }, onBack = { navController.popBackStack() }) }
        composable<AppRoute.AccountChoice> { AccountChoiceScreen(onPersonal = { navController.navigate(AppRoute.PlaceSetup) }, onOrganization = { navController.navigate(AppRoute.PlaceSetup) }, onInvitation = { navController.navigate(AppRoute.PlaceSetup) }) }
        composable<AppRoute.PlaceSetup> { PlaceSetupScreen(onContinue = { navController.navigate(AppRoute.Interests) }) }
        composable<AppRoute.Interests> { InterestsScreen(onContinue = { navController.navigate(AppRoute.SetupComplete) }, onSkip = { navController.navigate(AppRoute.SetupComplete) }) }
        composable<AppRoute.SetupComplete> { SetupCompleteScreen(onFinish = { navController.navigate(AppRoute.Today) { popUpTo(AppRoute.Welcome) { inclusive = true } } }) }

        // Main App
        composable<AppRoute.Today> { MainAppScaffold(currentRoute = "today", onNavigate = { navController.navigate(it) }) { TodayRoute() } }
        composable<AppRoute.Explore> { MainAppScaffold(currentRoute = "explore", onNavigate = { navController.navigate(it) }) { ExploreScreen() } }
        composable<AppRoute.Updates> { MainAppScaffold(currentRoute = "updates", onNavigate = { navController.navigate(it) }) { UpdatesScreen() } }
        
        composable<AppRoute.Report> { 
            MainAppScaffold(currentRoute = "report", onNavigate = { navController.navigate(it) }) { 
                ReportScreen(onSelectType = { type -> navController.navigate(AppRoute.ReportForm(type)) }) 
            } 
        }
        
        composable<AppRoute.Profile> { 
            MainAppScaffold(currentRoute = "profile", onNavigate = { navController.navigate(it) }) { 
                ProfileScreen(onOpenPlaces = { navController.navigate(AppRoute.SavedPlaces) }) 
            } 
        }

        // Drill-downs
        composable<AppRoute.ReportForm> { backStackEntry ->
            val route = backStackEntry.toRoute<AppRoute.ReportForm>()
            ReportFormScreen(reportType = route.reportType, onBack = { navController.popBackStack() }, onSubmitSuccess = { navController.navigate(AppRoute.ReportSuccess) })
        }
        composable<AppRoute.ReportSuccess> { 
            ReportSuccessScreen(onDone = { navController.popBackStack(route = AppRoute.Report, inclusive = false) }) 
        }
        composable<AppRoute.UpdateDetail> { UpdateDetailScreen(onBack = { navController.popBackStack() }) }
        composable<AppRoute.SavedPlaces> { SavedPlacesScreen(onBack = { navController.popBackStack() }) }
    }
}
