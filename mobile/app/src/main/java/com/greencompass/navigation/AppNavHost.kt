package com.greencompass.navigation

import androidx.compose.runtime.Composable
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.rememberNavController
import com.greencompass.feature.onboarding.LanguageSelectionScreen
import com.greencompass.feature.onboarding.WelcomeScreen

@Composable
fun AppNavHost() {
    val navController = rememberNavController()

    NavHost(
        navController = navController,
        startDestination = AppRoute.Welcome
    ) {
        composable<AppRoute.Welcome> {
            WelcomeScreen(
                onGetStarted = {
                    navController.navigate(AppRoute.AccountChoice)
                },
                onChooseLanguage = {
                    navController.navigate(AppRoute.LanguageSelection)
                }
            )
        }

        composable<AppRoute.LanguageSelection> {
            LanguageSelectionScreen(
                onContinue = {
                    navController.navigate(AppRoute.AccountChoice)
                },
                onBack = {
                    navController.popBackStack()
                }
            )
        }
        
        composable<AppRoute.AccountChoice> {
            // Placeholder for Account Choice screen
        }
    }
}
