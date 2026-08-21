package com.greencompass.feature.onboarding

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.greencompass.core.ui.GreenCompassColors

@Composable
fun WelcomeScreen(
    onGetStarted: () -> Unit,
    onChooseLanguage: () -> Unit
) {
    Surface(
        modifier = Modifier.fillMaxSize(),
        color = GreenCompassColors.WarmWhite
    ) {
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(32.dp),
            horizontalAlignment = Alignment.CenterHorizontally,
            verticalArrangement = Arrangement.Center
        ) {
            Box(
                modifier = Modifier
                    .size(80.dp)
                    .padding(bottom = 32.dp),
                contentAlignment = Alignment.Center
            ) {
                Text(text = "🌿", fontSize = 48.sp)
            }

            Text(
                text = "Green Compass",
                fontSize = 24.sp,
                fontWeight = FontWeight.SemiBold,
                color = GreenCompassColors.DeepForest,
                modifier = Modifier.padding(bottom = 24.dp)
            )

            Text(
                text = "Know Your Place.\nMove with Change.",
                fontSize = 32.sp,
                fontWeight = FontWeight.Bold,
                color = GreenCompassColors.Charcoal,
                textAlign = TextAlign.Center,
                lineHeight = 40.sp,
                modifier = Modifier.padding(bottom = 16.dp)
            )

            Text(
                text = "Clear environmental updates for the places\nthat matter to you.",
                fontSize = 16.sp,
                color = GreenCompassColors.MutedText,
                textAlign = TextAlign.Center,
                lineHeight = 24.sp,
                modifier = Modifier.padding(bottom = 48.dp)
            )

            Spacer(modifier = Modifier.weight(1f))

            Button(
                onClick = onGetStarted,
                modifier = Modifier
                    .fillMaxWidth()
                    .height(56.dp),
                colors = ButtonDefaults.buttonColors(
                    containerColor = GreenCompassColors.ForestGreen,
                    contentColor = Color.White
                ),
                shape = RoundedCornerShape(12.dp)
            ) {
                Text(
                    text = "Get started",
                    fontSize = 16.sp,
                    fontWeight = FontWeight.Medium
                )
            }

            Spacer(modifier = Modifier.height(16.dp))

            TextButton(
                onClick = onChooseLanguage,
                modifier = Modifier.fillMaxWidth()
            ) {
                Text(
                    text = "Choose language",
                    fontSize = 15.sp,
                    color = GreenCompassColors.ForestGreen,
                    fontWeight = FontWeight.Medium
                )
            }
            
            Spacer(modifier = Modifier.height(24.dp))
        }
    }
}
