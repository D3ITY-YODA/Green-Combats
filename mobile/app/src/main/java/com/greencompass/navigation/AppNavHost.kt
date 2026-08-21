package com.greencompass.navigation

import androidx.compose.runtime.Composable
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.rememberNavController
import com.greencompass.feature.onboarding.*
import com.greencompass.feature.today.TodayRoute

@Composable
fun AppNavHost() {
    val navController = rememberNavController()

    NavHost(navController = navController, startDestination = AppRoute.Welcome) {
        composable<AppRoute.Welcome> {
            WelcomeScreen(
                onGetStarted = { navController.navigate(AppRoute.AccountChoice) },
                onChooseLanguage = { navController.navigate(AppRoute.LanguageSelection) }
            )
        }
        composable<AppRoute.LanguageSelection> {
            LanguageSelectionScreen(
                onContinue = { navController.navigate(AppRoute.AccountChoice) },
                onBack = { navController.popBackStack() }
            )
        }
        composable<AppRoute.AccountChoice> {
            AccountChoiceScreen(
                onPersonal = { navController.navigate(AppRoute.PersonalRegistration) },
                onOrganization = { navController.navigate(AppRoute.OrganizationSearch) },
                onInvitation = { navController.navigate(AppRoute.InvitationAcceptance) },
                onSignIn = { navController.navigate(AppRoute.PhoneEmailSignIn) }
            )
        }
        composable<AppRoute.PersonalRegistration> {
            PersonalRegistrationScreen(
                onBack = { navController.popBackStack() },
                onContinueWithGoogle = { navController.navigate(AppRoute.GoogleSignIn) },
                onContinue = { navController.navigate(AppRoute.VerificationCode) }
            )
        }
        composable<AppRoute.GoogleSignIn> {
            GoogleSignInScreen(
                onBack = { navController.popBackStack() },
                onSuccess = { navController.navigate(AppRoute.VerificationCode) },
                onUsePhoneEmail = { 
                    navController.popBackStack() 
                    navController.navigate(AppRoute.PhoneEmailSignIn) 
                }
            )
        }
        composable<AppRoute.PhoneEmailSignIn> {
            PhoneEmailSignInScreen(
                onBack = { navController.popBackStack() },
                onContinue = { navController.navigate(AppRoute.VerificationCode) },
                onCreateAccount = { 
                    navController.popBackStack() 
                    navController.navigate(AppRoute.PersonalRegistration) 
                }
            )
        }
        composable<AppRoute.VerificationCode> {
            VerificationCodeScreen(
                onBack = { navController.popBackStack() },
                onVerify = { navController.navigate(AppRoute.LocationSelection) },
                onSendAgain = { }
            )
        }
        composable<AppRoute.OrganizationSearch> {
            OrganizationSearchScreen(
                onBack = { navController.popBackStack() },
                onSelectOrganization = { navController.navigate(AppRoute.OrganizationSelection) },
                onRequestAccess = { navController.navigate(AppRoute.RequestAccess) }
            )
        }
        composable<AppRoute.OrganizationSelection> {
            OrganizationSelectionScreen(
                onBack = { navController.popBackStack() },
                onContinue = { navController.navigate(AppRoute.AccessPending) }
            )
        }
        composable<AppRoute.RequestAccess> {
            RequestAccessScreen(
                onBack = { navController.popBackStack() },
                onRequestSent = { navController.navigate(AppRoute.AccessPending) }
            )
        }
        composable<AppRoute.InvitationAcceptance> {
            InvitationAcceptanceScreen(
                onBack = { navController.popBackStack() },
                onAccept = { navController.navigate(AppRoute.LocationSelection) },
                onDecline = { navController.popBackStack() }
            )
        }
        composable<AppRoute.AccessPending> {
            AccessPendingScreen(
                onBack = { navController.popBackStack() },
                onViewRequest = { },
                onReturnToToday = { navController.navigate(AppRoute.Today) }
            )
        }
        composable<AppRoute.LocationSelection> {
            LocationSelectionScreen(
                onBack = { navController.popBackStack() },
                onUseLocation = { navController.navigate(AppRoute.Interests) },
                onSearchPlace = { navController.navigate(AppRoute.SearchPlace) },
                onChooseOnMap = { navController.navigate(AppRoute.MapPlaceSelection) },
                onContinue = { navController.navigate(AppRoute.Interests) }
            )
        }
        composable<AppRoute.SearchPlace> {
            SearchPlaceScreen(
                onBack = { navController.popBackStack() },
                onSelectPlace = { navController.popBackStack(); navController.navigate(AppRoute.Interests) }
            )
        }
        composable<AppRoute.MapPlaceSelection> {
            MapPlaceSelectionScreen(
                onBack = { navController.popBackStack() },
                onConfirmPlace = { navController.popBackStack(); navController.navigate(AppRoute.Interests) }
            )
        }
        composable<AppRoute.Interests> {
            InterestsScreen(
                onBack = { navController.popBackStack() },
                onContinue = { navController.navigate(AppRoute.NotificationPreferences) },
                onSkip = { navController.navigate(AppRoute.NotificationPreferences) }
            )
        }
        composable<AppRoute.NotificationPreferences> {
            NotificationPreferencesScreen(
                onBack = { navController.popBackStack() },
                onContinue = { navController.navigate(AppRoute.PrivacyPermission) }
            )
        }
        composable<AppRoute.PrivacyPermission> {
            PrivacyPermissionScreen(
                onBack = { navController.popBackStack() },
                onAllowLocation = { navController.navigate(AppRoute.SetupComplete) },
                onChooseManually = { navController.navigate(AppRoute.SetupComplete) },
                onReadPrivacy = { }
            )
        }
        composable<AppRoute.SetupComplete> {
            SetupCompleteScreen(
                onFinish = { 
                    navController.navigate(AppRoute.Today) { 
                        popUpTo(AppRoute.Welcome) { inclusive = true } 
                    } 
                }
            )
        }
        composable<AppRoute.Today> {
            TodayRoute()
        }
    }
}
